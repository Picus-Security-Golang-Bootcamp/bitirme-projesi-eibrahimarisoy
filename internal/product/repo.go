package product

import (
	"fmt"
	"patika-ecommerce/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func (r *ProductRepository) Migration() {
	r.db.AutoMigrate(&model.Role{})
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// InsertProduct insert product
func (r *ProductRepository) InsertProduct(product *model.Product) error {
	result := r.db.Where("sku = ?", product.SKU).FirstOrCreate(product)
	if result.Error != nil {
		return result.Error
	}
	// fmt.Println(r.db.Model(product).Association("Categories").Count())
	// if product.CategoriesID != nil {
	// 	for _, categoryID := range product.CategoriesID {
	// 		fmt.Println("categoryID: ", categoryID)
	// 		r.db.Debug().Model(*product).Omit("Languages").Association("Categories").Append(&model.Category{Base: model.Base{ID: categoryID}})
	// 	}
	// }
	return nil

}

// GetProducts get all products
func (r *ProductRepository) GetProducts() (*[]model.Product, error) {
	products := new([]model.Product)
	result := r.db.Preload("Categories").Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("products: ", products)
	return products, nil
}
