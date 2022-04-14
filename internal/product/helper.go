package product

import (
	"fmt"

	"gorm.io/gorm"
)

// Search adds where to search keywords
func Search(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		fmt.Println("Search: ", search)
		if search != "" {
			db = db.Where("name ILIKE ?", "%"+search+"%")
			db = db.Or("sku ILIKE ?", "%"+search+"%")
		}
		return db
	}
}
