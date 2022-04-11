package order

import (
	"patika-ecommerce/internal/model"

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

	order := &model.Order{
		UserID: cart.UserID,
		CartID: cart.ID,
		Status: model.OrderStatusCompleted,
	}

	if err := r.db.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, item := range cart.Items {
		orderItem := &model.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Price:     item.Price,
		}

		if err := r.db.Create(orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return order, nil
}
