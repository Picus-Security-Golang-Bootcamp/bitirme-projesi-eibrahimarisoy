package model

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`

	Products []Product `json:"products" gorm:"many2many:product_categories;"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.Slug = slug.Make(c.Name)
	return nil
}
