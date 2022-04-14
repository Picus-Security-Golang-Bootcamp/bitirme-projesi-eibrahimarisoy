package model

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        *string `json:"name"`
	Slug        string  `json:"slug" gorm:"unique"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"type:decimal(20,2)"`
	Stock       *int64  `json:"stock"`
	SKU         *string `json:"sku" gorm:"unique"`

	Categories   []Category    `json:"categories" gorm:"many2many:product_categories; constraint:OnDelete:CASCADE"`
	CategoriesID []strfmt.UUID `json:"categories_id" gorm:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.Slug = slug.Make(*p.Name + "-" + *p.SKU)
	return nil
}

// ToString converts the product to string
func (p *Product) ToString() string {
	return "Product: " +
		"ID: " + p.ID.String() +
		"Name: " + *p.Name +
		"Description: " + p.Description +
		// "Price: " + fmt.Sprintf("%f", *p.Price) +
		"Categories: " + fmt.Sprintf("%v", p.Categories)
	// "Stock: " + strconv.Itoa(p.Stock)
	// "SKU: " + *p.SKU
}
