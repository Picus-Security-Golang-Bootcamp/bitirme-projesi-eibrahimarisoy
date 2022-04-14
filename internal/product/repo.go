package product

import (
	"fmt"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *gorm.DB
}

func (r *ProductRepository) Migration() {
	r.db.AutoMigrate(&model.Product{})

}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// InsertProduct insert product
func (r *ProductRepository) InsertProduct(product *model.Product) error {
	result := r.db.Create(&product)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetProducts get all products
func (r *ProductRepository) GetProducts(pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	var products []model.Product
	var totalRows int64

	query := r.db.Model(&model.Product{}).Scopes(Search(pagination.Q)).Count(&totalRows).Preload("Categories")
	query.Scopes(paginationHelper.Paginate(totalRows, pagination, r.db)).Find(&products)

	pagination.Rows = ProductsToResponse(&products)

	return pagination, nil
}

// GetProduct get a single product
func (r *ProductRepository) GetProduct(id uuid.UUID) (*model.Product, error) {
	product := new(model.Product)
	result := r.db.Preload("Categories").Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

// GetProductWithoutCategories get a single product
func (r *ProductRepository) GetProductWithoutCategories(id uuid.UUID) (*model.Product, error) {
	product := new(model.Product)
	result := r.db.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

// DeleteProduct delete a single product
func (r *ProductRepository) DeleteProduct(product *model.Product) error {
	// r.db.Model(&product).Association("Categories").Delete(&product)
	// r.db.Model(&product).Association("Categories").Delete(&product)

	result := r.db.Select(clause.Associations).Delete(&product)
	// result := r.db.Where(model.Product{}).Delete(&product)

	if result.Error != nil {
		return result.Error
	}
	fmt.Println("product: ", product)
	return nil
}

// UpdateProduct update a single product
func (r *ProductRepository) UpdateProduct(product *model.Product) error {
	tx := r.db.Begin()
	exProduct := new(model.Product)

	// get product
	err := tx.Where("id = ?", product.ID).Preload("Categories").First(&exProduct)

	if err.Error != nil {
		return err.Error
	}

	// delete all associated categories
	if err := tx.Model(&exProduct).Association("Categories").Delete(&exProduct.Categories); err != nil {
		return err
	}

	result := tx.Model(&product).Updates(&product)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}
