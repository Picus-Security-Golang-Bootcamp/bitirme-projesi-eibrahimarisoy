package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCanceled  OrderStatus = "cancelled"
)

type Order struct {
	Base

	UserID uuid.UUID `json:"user_id"`
	User   User      `json:"user"`

	Status OrderStatus `json:"status"`

	CartID uuid.UUID `json:"cart_id"`
	Cart   Cart      `json:"cart"`

	TotalPrice float64 `json:"total_price" gorm:"type:numeric(10,2)"`

	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	Base
	OrderID uuid.UUID `json:"order_id"`
	Order   Order     `json:"order"`

	ProductID uuid.UUID `json:"product_id"`
	Product   Product   `json:"product"`

	Price float64 `json:"price" gorm:"type:numeric(10,2)"`
}

func (o *Order) IsCancelable() bool {
	today := time.Now()
	lastDay := today.AddDate(0, 0, 14)
	if o.CreatedAt.Before(lastDay) && o.Status == OrderStatusCompleted {
		return true
	}
	return false
}
