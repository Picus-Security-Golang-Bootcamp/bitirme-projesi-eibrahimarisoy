package order

import (
	"fmt"
	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepositoryInterface interface {
	CompleteOrder(user *model.User, cartId uuid.UUID) (*model.Order, error)
	GetOrdersByUser(user *model.User, pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error)
	GetOrderByIdAndUser(user *model.User, id uuid.UUID) (*model.Order, error)
	CancelOrder(id uuid.UUID, user *model.User) error
}

type OrderRepository struct {
	db *gorm.DB
}

func (r *OrderRepository) Migration() {
	r.db.AutoMigrate(&model.Order{})
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type OrderItemRepository struct {
	db *gorm.DB
}

func (r *OrderItemRepository) Migration() {
	r.db.AutoMigrate(&model.OrderItem{})
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

// CompleteOrder
func (r *OrderRepository) CompleteOrder(user *model.User, cartId uuid.UUID) (*model.Order, error) {
	zap.L().Debug("order.repo.CompleteOrder", zap.Reflect("user", user), zap.Reflect("cartId", cartId))

	tx := r.db.Begin()
	cart := model.Cart{}

	// get cart by id and user id
	if err := tx.Model(model.Cart{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Items.Product").
		Where("id = ? AND user_id = ? AND status = ?", cartId, user.ID, model.CartStatusCreated).
		Find(&cart).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// create order from cart
	order := model.Order{
		UserID:     cart.UserID,
		CartID:     cart.ID,
		Status:     model.OrderStatusCompleted,
		TotalPrice: cart.GetTotalPrice(),
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// create order items from cart items
	for _, item := range cart.Items {
		product := item.Product
		if *product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("product %s stock is not enough", *item.Product.Name)
		}
		newProd := &model.Product{}
		// update product stock
		if err := tx.Model(&model.Product{}).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", item.Product.ID).
			First(&newProd).
			Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		// create order item
		for i := 0; i < int(item.Quantity); i++ {
			orderItem := &model.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Price:     item.Price,
			}

			if err := tx.Create(orderItem).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	// save cart
	cart.Status = model.CartStatusPaid
	if err := tx.Save(cart).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &order, nil
}

// GetOrdersByUser returns all orders of a user
func (r *OrderRepository) GetOrdersByUser(user *model.User, pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	zap.L().Debug("order.repo.GetOrdersByUser", zap.Reflect("user", user), zap.Reflect("pagination", pagination))

	var (
		orders    []*model.Order
		totalRows int64
	)

	query := r.db.Model(&model.Order{}).Where("user_id = ?", user.ID).Count(&totalRows).Preload("Items.Product")
	query.Scopes(paginationHelper.Paginate(totalRows, pagination, r.db)).Find(&orders)
	pagination.Rows = OrdersToOrderDetailedResponse(orders)

	return pagination, nil
}

// GetOrderByIdAndUser returns an order by id and user
func (r *OrderRepository) GetOrderByIdAndUser(user *model.User, id uuid.UUID) (*model.Order, error) {
	zap.L().Debug("order.repo.GetOrderByIdAndUser", zap.Reflect("user", user), zap.Reflect("id", id))

	var order model.Order
	if err := r.db.Preload("Items.Product").Where("id = ? AND user_id = ?", id, user.ID).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// CancelOrder cancels an order
func (r *OrderRepository) CancelOrder(id uuid.UUID, user *model.User) error {
	zap.L().Debug("order.repo.CancelOrder", zap.Reflect("user", user), zap.Reflect("id", id))

	tx := r.db.Begin()
	var order model.Order
	// get order by id and user id
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Items.Product").
		Where("id = ? AND user_id = ? AND status = ?", id, user.ID, model.OrderStatusCompleted).
		First(&order).Error; err != nil {
		tx.Rollback()
		return err
	}
	// check if order is cancelable
	if !order.IsCancelable() {
		return httpErr.OrderCannotBeCanceledError
	}

	// update products stock
	newProd := &model.Product{}
	for _, item := range order.Items {
		if err := tx.Model(&model.Product{}).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", item.Product.ID).
			First(&newProd).
			Update("stock", gorm.Expr("stock + ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// update order status
	order.Status = model.OrderStatusCanceled
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
