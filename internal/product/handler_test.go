package product

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func getProductPOSTPayload() []byte {
	var jsonStr = []byte(
		`{
			"name": "product name",
			"description": "product description",
			"price": 1450.444,
			"sku": "PRODUCT-SKU",
			"stock": 40,
			"categories": []
		}`)

	return jsonStr
}

func getProductPUTPayload() []byte {
	var jsonStr = []byte(
		`{
			"name": "product name update",
			"description": "product description update",
			"price": 1450.444,
			"sku": "PRODUCT-SKU-UPDATE",
			"stock": 40,
			"categories": []
		}`)

	return jsonStr
}

func Test_productHandler_createProduct(t *testing.T) {

	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"

	id := uuid.New()

	gin.SetMode(gin.TestMode)

	categoryRepo := &CategoryMockRepository{
		Items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &categoryName,
				Description: categoryDescription,
			},
		},
	}

	mockProductRepo := &ProductMockRepository{
		items: []model.Product{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}

	w := httptest.NewRecorder()
	productHandler := &productHandler{
		productRepo:  mockProductRepo,
		categoryRepo: categoryRepo,
	}
	c, r := gin.CreateTestContext(w)
	r.POST("/products", productHandler.createProduct)
	c.Request, _ = http.NewRequest("GET", "/products", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	productHandler.createProduct(c)

	assert.Equal(t, http.StatusCreated, w.Code)

}

type ProductMockRepository struct {
	items      []model.Product
	categories []model.Category
}

func (r *ProductMockRepository) Insert(product *model.Product) error {
	fmt.Println("InsertProduct: ", product)

	for _, item := range product.Categories {
		for _, category = range r.category.Items {
			if item.ID == category.ID {
				return errors.New("category already exists")
			}
		}

	}

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
func (r *ProductMockRepository) GetAll(pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	var products []model.Product
	var totalRows int64

	query := r.db.Model(&model.Product{}).Scopes(Search(pagination.Q)).Count(&totalRows).Preload("Categories")
	query.Scopes(paginationHelper.Paginate(totalRows, pagination, r.db)).Find(&products)

	pagination.Rows = ProductsToResponse(&products)

	return pagination, nil
}

// GetProduct get a single product
func (r *ProductMockRepository) Get(id uuid.UUID) (*model.Product, error) {
	product := new(model.Product)
	result := r.db.Preload("Categories").Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

// GetProductWithoutCategories get a single product
func (r *ProductMockRepository) GetProductWithoutCategories(id uuid.UUID) (*model.Product, error) {
	product := new(model.Product)
	result := r.db.Where("id = ?", id).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

// DeleteProduct delete a single product
func (r *ProductMockRepository) Delete(product model.Product) error {
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
func (r *ProductMockRepository) Update(product *model.Product) error {
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

///Categories
type CategoryMockRepository struct {
	Items []model.Category
}

func (r *CategoryMockRepository) InsertCategory(category *model.Category) error {
	for _, item := range r.Items {
		item.ID = category.ID
		return errors.New("category already exists")
	}
	r.Items = append(r.Items, *category)

	return nil
}

// GetCategories returns all categories
func (r *CategoryMockRepository) GetCategories() (*[]model.Category, error) {
	return &r.Items, nil
}

// GetCategoryByID returns a category by id
func (r *CategoryMockRepository) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	category := &model.Category{}
	for _, item := range r.Items {
		if item.ID == id {
			category = &item
		}
	}
	return category, nil
}

// UpdateCategory updates a category with the given id
func (r *CategoryMockRepository) UpdateCategory(category *model.Category) error {
	for i, item := range r.Items {
		if item.ID == category.ID {
			r.Items[i] = *category
		}
	}

	return nil
}

// InsertBulkCategory inserts bulk categories //TODO
func (r *CategoryMockRepository) InsertBulkCategory(categories *[]model.Category) error {
	for _, category := range *categories {
		for _, item := range r.Items {
			if item.Name == category.Name {
				return errors.New("category already exists")
			}
		}
		r.Items = append(r.Items, category)
	}

	return nil
}

// DeleteCategory deletes a category by id
func (r *CategoryMockRepository) Delete(category *model.Category) error {
	isExist := false
	for index, item := range r.Items {
		if item.ID == category.ID {
			category = &item
			isExist = true
			r.Items = append(r.Items[:index], r.Items[index+1:]...)
		}
	}
	if !isExist {
		return errors.New("category not found")
	}
	return nil
}
