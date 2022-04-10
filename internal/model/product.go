package model

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        *string  `json:"name"`
	Slug        string   `json:"slug" gorm:"unique"`
	Description string   `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int64   `json:"stock"`
	SKU         *string  `json:"sku" gorm:"unique"`

	Categories   *[]Category   `json:"categories" gorm:"many2many:product_categories"`
	CategoriesID []strfmt.UUID `json:"categories_id" gorm:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.Slug = slug.Make(*p.Name)
	return nil
}

// AfterDelete hook defined for cascade delete
func (p *Product) AfterDelete(tx *gorm.DB) error {
	return tx.Model("product_categories").Where("product_id = ?", p.ID).Unscoped().Delete(&p).Error
}

// ToString converts the product to string
func (p *Product) ToString() string {
	return "Product: " +
		"ID: " + p.ID.String() +
		"Name: " + *p.Name +
		"Description: " + p.Description +
		"Price: " + fmt.Sprintf("%f", *p.Price)
	// "Stock: " + strconv.Itoa(p.Stock)
	// "SKU: " + *p.SKU
}
