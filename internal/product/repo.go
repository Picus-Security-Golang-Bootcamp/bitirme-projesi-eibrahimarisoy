package product

import (
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"

	"github.com/google/uuid"
	"go.uber.org/zap"
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
	zap.L().Debug("product.repo.Insert", zap.Reflect("product", product))

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
	zap.L().Debug("product.repo.GetAll", zap.Reflect("pagination", pagination))

	var products []model.Product
	var totalRows int64

	query := r.db.Model(&model.Product{}).Scopes(Search(pagination.Q)).Count(&totalRows).Preload("Categories")
	query.Scopes(paginationHelper.Paginate(totalRows, pagination, r.db)).Find(&products)

	pagination.Rows = ProductsToResponse(&products)

	return pagination, nil
}

// GetProduct get a single product
func (r *ProductRepository) Get(id uuid.UUID) (*model.Product, error) {
	zap.L().Debug("product.repo.Get", zap.Reflect("id", id))

	product := new(model.Product)
	result := r.db.Preload("Categories").Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

// GetProductWithoutCategories get a single product
func (r *ProductRepository) GetProductWithoutCategories(id uuid.UUID) (*model.Product, error) {
	zap.L().Debug("product.repo.GetProductWithoutCategories", zap.Reflect("id", id))

	product := new(model.Product)
	result := r.db.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

// DeleteProduct delete a single product
func (r *ProductRepository) Delete(product *model.Product) error {
	zap.L().Debug("product.repo.Delete", zap.Reflect("product", product))

	result := r.db.Select(clause.Associations).Delete(product)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateProduct update a single product
func (r *ProductRepository) Update(product *model.Product) error {
	zap.L().Debug("product.repo.Update", zap.Reflect("product", product))

	tx := r.db.Begin()
	for index, item := range product.Categories {
		category := new(model.Category)
		if err := tx.Model(category).Where("id = ?", item.ID).First(&category); err != nil {
			tx.Rollback()
			return err.Error
		}
		product.Categories[index] = *category
	}

	exProduct := new(model.Product)

	// get ex product
	if err := tx.Where("id = ?", product.ID).Preload("Categories").First(&exProduct); err.Error != nil {
		tx.Rollback()
		return err.Error
	}

	// delete all associated categories
	if err := tx.Model(&exProduct).Association("Categories").Delete(&exProduct.Categories); err != nil {
		tx.Rollback()
		return err
	}

	if result := tx.Model(&product).Updates(&product); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// update cart items price if product price is changed
	if exProduct.Price != product.Price {
		if err := tx.Model(&model.CartItem{}).Where("product_id = ?", product.ID).Update("price", product.Price).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
