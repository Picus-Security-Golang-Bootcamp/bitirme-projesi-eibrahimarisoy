package product

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

const (
	categoryId = "85410bda-d3eb-475c-9d4b-f6d1ee9e4b7f"
)

func getProductPOSTPayload() []byte {
	var jsonStr = []byte(
		`{
			"name": "product name",
			"description": "product description",
			"price": 1450.444,
			"sku": "PRODUCT-SKU",
			"stock": 40,
			"categories": [
				{"id": "` + categoryId + `"}
			]
		}`)

	return jsonStr
}

func getProductPUTPayload() []byte {
	var jsonStr = []byte(
		`{
			"name": "product name update",
			"description": "product description update",
			"price": 1450.44,
			"sku": "PRODUCT-SKU-UPDATE",
			"stock": 40,
			"categories": []
		}`)

	return jsonStr
}

func Test_productHandler_createProduct(t *testing.T) {
	// name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	// id := uuid.New()
	gin.SetMode(gin.TestMode)

	catId, _ := uuid.Parse(fmt.Sprintf("%s", categoryId))
	mockProductRepo := &ProductMockRepository{
		items: []model.Product{
			// {
			// 	Base:        model.Base{ID: id},
			// 	Name:        &name,
			// 	Description: description,
			// 	Price:       0,
			// 	Stock:       new(int64),
			// 	SKU:         new(string),
			// 	Categories: []model.Category{
			// 		{
			// 			Base:        model.Base{ID: catId},
			// 			Name:        &categoryName,
			// 			Description: categoryDescription,
			// 		},
			// 	},
			// },
		},
		categories: []model.Category{
			{
				Base:        model.Base{ID: catId},
				Name:        &categoryName,
				Description: categoryDescription,
			},
		},
	}

	w := httptest.NewRecorder()
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	c, r := gin.CreateTestContext(w)
	r.POST("/products", productHandler.createProduct)
	c.Request, _ = http.NewRequest("GET", "/products", nil)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPOSTPayload()))
	c.Request.Header.Set("Content-Type", "application/json")
	productHandler.createProduct(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, len(mockProductRepo.items))
}

func Test_productHandler_getProducts(t *testing.T) {
	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	pagination := paginationHelper.Pagination{
		Limit: 2,
		Page:  1,
	}

	id := uuid.New()
	gin.SetMode(gin.TestMode)

	catId, _ := uuid.Parse(fmt.Sprintf("%s", categoryId))
	mockProductRepo := &ProductMockRepository{
		items: []model.Product{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
				Price:       0,
				Stock:       new(int64),
				SKU:         new(string),
				Categories: []model.Category{
					{
						Base:        model.Base{ID: catId},
						Name:        &categoryName,
						Description: categoryDescription,
					},
				},
			},
		},
		categories: []model.Category{},
	}

	w := httptest.NewRecorder()
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	c, r := gin.CreateTestContext(w)
	c.Set("pagination", &pagination)
	r.GET("/products", productHandler.getProducts)
	c.Request, _ = http.NewRequest("GET", "/products", nil)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPOSTPayload()))
	c.Request.Header.Set("Content-Type", "application/json")
	productHandler.getProducts(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_productHandler_getProduct(t *testing.T) {
	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	id := uuid.New()
	catId := uuid.New()
	mockProductRepo := &ProductMockRepository{
		items: []model.Product{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
				Price:       0,
				Stock:       new(int64),
				SKU:         new(string),
				Categories: []model.Category{
					{
						Base:        model.Base{ID: catId},
						Name:        &categoryName,
						Description: categoryDescription,
					},
				},
			},
		},
		categories: []model.Category{},
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	c, r := gin.CreateTestContext(w)
	r.GET("/product/:id", productHandler.getProduct)
	c.Params = []gin.Param{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("GET", "/product/:id", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	productHandler.getProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_productHandler_deleteProduct(t *testing.T) {
	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	id := uuid.New()
	catId := uuid.New()
	mockProductRepo := &ProductMockRepository{
		items: []model.Product{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
				Price:       0,
				Stock:       new(int64),
				SKU:         new(string),
				Categories: []model.Category{
					{
						Base:        model.Base{ID: catId},
						Name:        &categoryName,
						Description: categoryDescription,
					},
				},
			},
		},
		categories: []model.Category{},
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	c, r := gin.CreateTestContext(w)
	r.DELETE("/product/:id", productHandler.deleteProduct)
	c.Params = []gin.Param{{Key: "id", Value: id.String()}}
	c.Request, _ = http.NewRequest("DELETE", "/product/:id", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	productHandler.deleteProduct(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, 0, len(mockProductRepo.items))
}

func Test_productHandler_updateProduct(t *testing.T) {
	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	id := uuid.New()
	catId := uuid.New()
	mockProductRepo := &ProductMockRepository{
		items: []model.Product{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
				Price:       0,
				Stock:       new(int64),
				SKU:         new(string),
				Categories: []model.Category{
					{
						Base:        model.Base{ID: catId},
						Name:        &categoryName,
						Description: categoryDescription,
					},
				},
			},
		},
		categories: []model.Category{},
	}
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	c, r := gin.CreateTestContext(w)
	r.PUT("/product/:id", productHandler.updateProduct)
	c.Request, _ = http.NewRequest("PUT", "/product/:id", nil)
	c.Params = []gin.Param{{Key: "id", Value: id.String()}}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPUTPayload()))
	c.Request.Header.Set("Content-Type", "application/json")
	productHandler.updateProduct(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "product name update", mockProductRepo.items[0].Name)
}

type ProductMockRepository struct {
	items      []model.Product
	categories []model.Category
}

func (r *ProductMockRepository) Insert(product *model.Product) error {

	newCategories := []model.Category{}

	for _, item := range product.Categories {
		for _, category := range r.categories {
			if item.ID == category.ID {
				newCategories = append(newCategories, category)
			}
		}
	}

	if len(newCategories) != len(product.Categories) {
		return errors.New("category not found")
	}

	product.Categories = newCategories

	r.items = append(r.items, *product)

	return nil
}

// GetProducts get all products
func (r *ProductMockRepository) GetAll(pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	var products []model.Product
	pagination.TotalRows = int64(len(r.items))
	pagination.Rows = ProductsToResponse(&products)

	return pagination, nil
}

// GetProduct get a single product
func (r *ProductMockRepository) Get(id uuid.UUID) (*model.Product, error) {
	for _, item := range r.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, nil
}

// GetProductWithoutCategories get a single product
func (r *ProductMockRepository) GetProductWithoutCategories(id uuid.UUID) (*model.Product, error) {
	// product := new(model.Product)
	// result := r.db.Where("id = ?", id).First(&product)
	// if result.Error != nil {
	// 	return nil, result.Error
	// }

	return nil, nil
}

// DeleteProduct delete a single product
func (r *ProductMockRepository) Delete(product model.Product) error {
	for i, item := range r.items {
		if item.ID == product.ID {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("product not found")
}

// UpdateProduct update a single product
func (r *ProductMockRepository) Update(product *model.Product) error {
	for i, item := range r.items {
		if item.ID == product.ID {
			r.items[i] = *product
			return nil
		}
	}
	return errors.New("product not found")
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
