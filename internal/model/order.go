package model

import "github.com/go-openapi/strfmt"

type OrderStatus string

const (
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	Base

	UserID strfmt.UUID `json:"user_id"`
	User   User        `json:"user"`

	Status OrderStatus `json:"status"`

	CartID strfmt.UUID `json:"cart_id"`
	Cart   Cart        `json:"cart"`

	TotalPrice float64 `json:"total_price"`

	Items []OrderItem `json:"items"`
}

type OrderItem struct {
	Base
	OrderID strfmt.UUID `json:"order_id"`
	Order   Order       `json:"order"`

	ProductID strfmt.UUID `json:"product_id"`
	Product   Product     `json:"product"`

	Price float64 `json:"price"`
}

func (o *Order) IsCancelable() bool {
	// TODO TODO TODO
	return o.Status == OrderStatusCompleted
}
