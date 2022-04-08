package model

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`

	Categories []Category `json:"category" gorm:"many2many:product_categories;"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.Slug = slug.Make(p.Name)
	return nil
}
