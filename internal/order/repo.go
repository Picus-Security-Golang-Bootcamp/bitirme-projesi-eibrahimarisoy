package order

import (
	"errors"
	"fmt"
	"patika-ecommerce/internal/model"

	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"
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
func (r *OrderRepository) CompleteOrder(cart *model.Cart) (*model.Order, error) {
	tx := r.db.Begin() // TODO total price

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
			fmt.Println("PRODUCY^", item.Product)

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
func (r *OrderRepository) GetOrdersByUser(user *model.User) ([]*model.Order, error) {
	var orders []*model.Order
	if err := r.db.Preload("Items.Product").Where("user_id = ?", user.ID).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByIdAndUser returns an order by id and user
func (r *OrderRepository) GetOrderByIdAndUser(user *model.User, id strfmt.UUID) (*model.Order, error) {
	var order model.Order
	if err := r.db.Where("id = ? AND user_id = ?", id, user.ID).First(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// CancelOrder cancels an order
func (r *OrderRepository) CancelOrder(order *model.Order) error {
	// TODO product stock update
	if !order.IsCancelable() {
		// return model.ErrOrderCannotBeCanceled
		return errors.New("order cannot be canceled")
	}

	order.Status = model.OrderStatusCancelled
	if err := r.db.Save(order).Error; err != nil {
		return err
	}

	return nil
}
