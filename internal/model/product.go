package model

import (
	"github.com/go-openapi/strfmt"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        *string  `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Price       *float64 `json:"price"`
	Stock       *int     `json:"stock"`
	SKU         *string  `json:"sku " gorm:"unique"`

	Categories   *[]Category   `json:"categories" gorm:"many2many:product_categories;"`
	CategoriesID []strfmt.UUID `json:"categories_id" gorm:"-"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.Slug = slug.Make(*p.Name)
	return nil
}
