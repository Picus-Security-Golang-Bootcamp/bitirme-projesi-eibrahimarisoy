package product

import (
	"fmt"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MockProductRepository interface {
	Insert(product *model.Product) error
	GetAll(pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error)
	Get(id uuid.UUID) (*model.Product, error)
	GetProductWithoutCategories(id uuid.UUID) (*model.Product, error)
	Delete(product *model.Product) error
	Update(product *model.Product) error
}

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
func (r *ProductRepository) Insert(product *model.Product) error {
	fmt.Println("InsertProduct: ", product)
	tx := r.db.Begin()

	result := tx.Omit("Categories").Create(product)
	if err := result.Error; err != nil {
		tx.Rollback()
		return err
	}
	// insert categories
	for _, category := range product.Categories {
		if err := tx.Model(&category).Association("Products").Append(product); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// GetProducts get all products
func (r *ProductRepository) GetAll(pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	var products []model.Product
	var totalRows int64

	query := r.db.Model(&model.Product{}).Scopes(Search(pagination.Q)).Count(&totalRows).Preload("Categories")
	query.Scopes(paginationHelper.Paginate(totalRows, pagination, r.db)).Find(&products)

	pagination.Rows = ProductsToResponse(&products)

	return pagination, nil
}

// GetProduct get a single product
func (r *ProductRepository) Get(id uuid.UUID) (*model.Product, error) {
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
func (r *ProductRepository) Delete(product *model.Product) error {
	// r.db.Model(&product).Association("Categories").Delete(&product)
	// r.db.Model(&product).Association("Categories").Delete(&product)

	result := r.db.Select(clause.Associations).Delete(product)
	// result := r.db.Where(model.Product{}).Delete(&product)

	if result.Error != nil {
		return result.Error
	}
	fmt.Println("product: ", product)
	return nil
}

// UpdateProduct update a single product
func (r *ProductRepository) Update(product *model.Product) error {
	tx := r.db.Begin()
	for index, item := range product.Categories {
		category := new(model.Category)
		result := tx.Model(category).Where("id = ?", item.ID).First(&category)
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
		product.Categories[index] = *category
	}

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
