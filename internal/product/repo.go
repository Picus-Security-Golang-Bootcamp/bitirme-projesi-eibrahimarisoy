package product

import (
	"errors"
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
	result := r.db.Where("sku = ?", product.SKU).Create(product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			return fmt.Errorf("product with sku %s already exists", *product.SKU)
		}
	}

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
