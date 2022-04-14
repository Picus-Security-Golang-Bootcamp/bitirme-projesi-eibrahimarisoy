package product

import (
	"fmt"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"

	"github.com/go-openapi/strfmt"
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
	fmt.Println("GetProducts: ", pagination)
	var products []*model.Product

	r.db.Scopes(paginationHelper.Paginate(products, pagination, r.db)).Find(&products)
	pagination.Rows = products

	return pagination, nil

	// products := new([]model.Product)

	// totalRows, totalPages := int64(0), 0

	// offset := (pagination.Page - 1) * pagination.Limit

	// // generate where query
	// search := pagination.Q

	// find := r.db
	// find = find.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Preload("Categories")
	// if search != "" {
	// 	find := find.Scopes(Search(search))
	// 	fmt.Println("search: ", find)
	// }

	// find = find.Find(products)
	// fmt.Println("find: ", products)

	// pagination.Rows = products
	// err := find.Count(&totalRows).Error

	// if err != nil {
	// 	return nil, 0, err
	// }

	// pagination.TotalRows = totalRows

	// // calculate total pages
	// totalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))

	// return pagination, totalPages, nil
}

// GetProduct get a single product
func (r *ProductRepository) GetProduct(id strfmt.UUID) (*model.Product, error) {
	product := new(model.Product)
	result := r.db.Preload("Categories").Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("product: ", product)
	return product, nil
}

// GetProductWithoutCategories get a single product
func (r *ProductRepository) GetProductWithoutCategories(id strfmt.UUID) (*model.Product, error) {
	product := new(model.Product)
	result := r.db.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println("product: ", product)
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
