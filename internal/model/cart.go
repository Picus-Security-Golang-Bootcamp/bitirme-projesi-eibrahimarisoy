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

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.Status = CartStatusCreated
	return nil
}

type CartItem struct {
	Base
	CartID strfmt.UUID `json:"cart_id"`
	Cart   Cart        `json:"cart"`

	ProductID strfmt.UUID `json:"product_id"`
	Product   Product     `json:"product"`

	Quantity int `json:"quantity"`
	Price    int `json:"price"`
}
