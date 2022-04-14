package model

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

	Quantity int64           `json:"quantity" gorm:"not null"`
	Price    decimal.Decimal `json:"price" gorm:"not null"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	c.Status = CartStatusCreated
	return nil
}

func (c *Cart) GetCartItemByID(id uuid.UUID) (*CartItem, error) {
	for _, item := range c.Items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("Cart item not found")
}

func (c *Cart) GetTotalPrice() decimal.Decimal {
	var totalPrice decimal.Decimal
	for _, item := range c.Items {
		totalPrice.Add(item.Price.Mul(decimal.NewFromInt(item.Quantity)))
	}
	return totalPrice
}
