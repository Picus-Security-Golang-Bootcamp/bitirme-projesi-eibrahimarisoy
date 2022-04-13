package model

import (
	"fmt"

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
	Status CartStatus `json:"status" gorm:"type:varchar(10);not null"`

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

	Quantity int64   `json:"quantity" gorm:"not null"`
	Price    float64 `json:"price" gorm:"not null"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.Status = CartStatusCreated
	return nil
}

func (c *Cart) GetCartItemByID(id strfmt.UUID) (*CartItem, error) {
	for _, item := range c.Items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("Cart item not found")
}

func (c *Cart) GetTotalPrice() float64 {
	var totalPrice float64
	for _, item := range c.Items {
		totalPrice += item.Price * float64(item.Quantity)
	}
	return totalPrice
}
