package order

import (
	"errors"
	"fmt"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
	tx := r.db.Begin() // TODO total price
	cart := model.Cart{}

	if err := tx.Model(model.Cart{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Items.Product").
		Where("id = ? AND user_id = ? AND status = ?", cartId, user.ID, model.CartStatusCreated).
		Find(&cart).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

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

	for _, item := range cart.Items {
		product := item.Product
		if *product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("product %s stock is not enough", item.Product.Name)
		}

		*product.Stock -= item.Quantity
		tx.Save(&product)

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
	var order model.Order
	if err := r.db.Preload("Items.Product").Where("id = ? AND user_id = ?", id, user.ID).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// CancelOrder cancels an order
func (r *OrderRepository) CancelOrder(id uuid.UUID, user *model.User) error {
	tx := r.db.Begin()
	var order model.Order

	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Items.Product").
		Where("id = ? AND user_id = ? AND status = ?", id, user.ID, model.OrderStatusCompleted).
		First(&order).Error; err != nil {

		tx.Rollback()
		return err
	}
	if !order.IsCancelable() {
		// return model.ErrOrderCannotBeCanceled
		return errors.New("order cannot be canceled")
	}

	for _, item := range order.Items {
		product := item.Product
		*product.Stock += 1
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	order.Status = model.OrderStatusCanceled
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
