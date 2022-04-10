package model

import (
	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"
)

type CartStatus string

const (
	CartStatusCreated   CartStatus = "created"
	CartStatusPaid      CartStatus = "paid"
	CartStatusCancelled CartStatus = "cancelled"
)

type Cart struct {
	Base
	Status CartStatus `json:"status"`

	UserID strfmt.UUID `json:"user_id"`
	User   User        `json:"user"`

	Items []CartItem `json:"items"`
}
type CartItem struct {
	Base
	CartID strfmt.UUID `json:"cart_id"`
	Cart   Cart        `json:"cart"`

	ProductID strfmt.UUID `json:"product_id"`
	Product   Product     `json:"product"`

	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.Status = CartStatusCreated
	return nil
}

func (c *Cart) CanAddProduct(id strfmt.UUID, quantity int) error {
	// TODO
	item := CartItem{
		CartID:    c.ID,
		ProductID: id,
		Quantity:  quantity,
	}
	c.Items = append(c.Items, item)

	return nil
}
