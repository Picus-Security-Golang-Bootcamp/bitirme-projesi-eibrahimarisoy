package model

type CartStatus string

const (
	CartStatusCreated   CartStatus = "created"
	CartStatusPaid      CartStatus = "paid"
	CartStatusCancelled CartStatus = "cancelled"
)

type Cart struct {
	Base
	Status CartStatus `json:"status"`

	UserID uint `json:"user_id"`
	User   User `json:"user"`
}

type CartItem struct {
	Base
	CartID uint `json:"cart_id"`
	Cart   Cart `json:"cart"`

	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`

	Quantity int `json:"quantity"`
	Price    int `json:"price"`
}
