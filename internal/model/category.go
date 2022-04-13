package model

type Category struct {
	Base
	Name        *string `json:"name" gorm:"type:varchar(100);not null;unique"`
	Slug        string  `json:"slug" gorm:"type:varchar(100);not null;unique"`
	Description string  `json:"description" gorm:"type:varchar(255)"`

	Products []Product `json:"products" gorm:"many2many:product_categories"`
}

// func (c *Category) BeforeCreate(tx *gorm.DB) error {
// 	if c.Slug == "" {
// 		c.Slug = slug.Make(*c.Name)
// 	}
// 	return nil
// }

// ToString converts the category to string
func (c *Category) ToString() string {
	return "Category: " +
		"ID: " + c.ID.String() +
		"Name: " + *c.Name +
		"Description: " + c.Description
}
