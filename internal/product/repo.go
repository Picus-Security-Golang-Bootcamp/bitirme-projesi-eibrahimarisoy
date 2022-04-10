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

// GetProduct get a single product
func (r *ProductRepository) GetProduct(sku string) (*model.Product, error) {
	product := new(model.Product)
	result := r.db.Preload("Categories").Where("sku = ?", sku).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("product: ", product)
	return product, nil
}

// DeleteProduct delete a single product
func (r *ProductRepository) DeleteProduct(sku string) error {
	product := new(model.Product)
	// r.db.Model(&product).Association("Categories").Clear()

	result := r.db.Omit("product_categories").Where("sku = ?", sku).Delete(&product)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println("product: ", product)
	return nil
}
