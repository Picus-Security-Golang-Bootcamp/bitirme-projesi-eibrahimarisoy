package model

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	Base

	UserID uint `json:"user_id"`
	User   User `json:"user"`

	Status OrderStatus `json:"status"`

	CartID uint `json:"cart_id"`
	Cart   Cart `json:"cart"`

	TotalPrice int `json:"total_price"`
}

type OrderItem struct {
	Base
	OrderID uint  `json:"order_id"`
	Order   Order `json:"order"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`

	Price int `json:"price"`
}
