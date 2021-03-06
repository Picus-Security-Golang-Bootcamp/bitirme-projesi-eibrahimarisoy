package model

import (
	"fmt"

	"github.com/google/uuid"
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

	UserID uuid.UUID `json:"user_id"`
	User   User      `json:"user"`

	Items []CartItem `json:"items"`
}
type CartItem struct {
	Base
	CartID uuid.UUID `json:"cart_id"`
	Cart   Cart      `json:"cart"`

	ProductID uuid.UUID `json:"product_id"`
	Product   Product   `json:"product"`

	Quantity int64   `json:"quantity" gorm:"not null"`
	Price    float64 `json:"price" gorm:"not null"`
}

// BeforeCreate hook
func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.Status = CartStatusCreated
	return nil
}

// GetCartItemByID returns cart item by id
func (c *Cart) GetCartItemByID(id uuid.UUID) (*CartItem, error) {
	for _, item := range c.Items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("Cart item not found")
}

// GetTotalPrice returns total price of cart
func (c *Cart) GetTotalPrice() float64 {
	var totalPrice float64
	for _, item := range c.Items {
		totalPrice += item.Price * float64(item.Quantity)
	}
	return totalPrice
}
