package model

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        *string `json:"name" gorm:"type:varchar(100);not null;unique"`
	Slug        string  `json:"slug" gorm:"type:varchar(100);not null;unique"`
	Description string  `json:"description" gorm:"type:varchar(255)"`

	Products []Product `json:"products" gorm:"many2many:product_categories"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.Slug == "" {
		c.Slug = slug.Make(*c.Name)
	}
	return nil
}
