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
