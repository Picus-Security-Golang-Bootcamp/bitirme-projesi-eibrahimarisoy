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
	"gorm.io/gorm"
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

func getProductPOSTPayloadNotValid() []byte {
	var jsonStr = []byte(
		`{
			"name_not": "product name update",
			"description_not": "product description update",
			"price": 1450.44,
			"sku": "PRODUCT-SKU-UPDATE",
			"stock": 40,
			"categories": []
		}`)

	return jsonStr
}

func Test_productHandler_createProduct(t *testing.T) {
	categoryName, categoryDescription := "test category", "test category description"

	t.Run("createProduct_Successful", func(t *testing.T) {
		catId, _ := uuid.Parse(categoryId)
		mockProductRepo := &mockProductRepository{
			items: []model.Product{},
			categories: []model.Category{
				{
					Base:        model.Base{ID: catId},
					Name:        &categoryName,
					Description: categoryDescription,
				},
			},
		}
		productHandler := &productHandler{
			productRepo: mockProductRepo,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/products", nil)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPOSTPayload()))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.createProduct(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, 1, len(mockProductRepo.items))
	})

	t.Run("createProduct_Failed_reqBody", func(t *testing.T) {
		catId, _ := uuid.Parse(categoryId)
		mockProductRepo := &mockProductRepository{
			items: []model.Product{},
			categories: []model.Category{
				{
					Base:        model.Base{ID: catId},
					Name:        &categoryName,
					Description: categoryDescription,
				},
			},
		}
		productHandler := &productHandler{
			productRepo: mockProductRepo,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/products", nil)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("failed-req-body")))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.createProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("createProduct_Failed_notValidate", func(t *testing.T) {
		catId, _ := uuid.Parse(categoryId)
		mockProductRepo := &mockProductRepository{
			items: []model.Product{},
			categories: []model.Category{
				{
					Base:        model.Base{ID: catId},
					Name:        &categoryName,
					Description: categoryDescription,
				},
			},
		}
		productHandler := &productHandler{
			productRepo: mockProductRepo,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/products", nil)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPOSTPayloadNotValid()))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.createProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("createProduct_Failed_notFoundCategory", func(t *testing.T) {
		mockProductRepo := &mockProductRepository{
			items:      []model.Product{},
			categories: []model.Category{},
		}
		productHandler := &productHandler{
			productRepo: mockProductRepo,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/products", nil)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPOSTPayload()))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.createProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func Test_productHandler_getProducts(t *testing.T) {
	id := uuid.New()
	catId := uuid.New()
	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	pagination := paginationHelper.Pagination{
		Limit: 2,
		Page:  1,
	}
	mockProductRepo := &mockProductRepository{
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
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	t.Run("getProducts_Successful", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Set("pagination", &pagination)
		c.Request, _ = http.NewRequest("GET", "/products", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.getProducts(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func Test_productHandler_getProduct(t *testing.T) {
	name, description := "test", "test"
	id, catId := uuid.New(), uuid.New()
	categoryName, categoryDescription := "test category", "test category description"

	mockProductRepo := &mockProductRepository{
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
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	t.Run("getProduct_Successful", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest("GET", "/product/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.getProduct(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("getProduct_Failed_idNotValid", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "id.String()"}}
		c.Request, _ = http.NewRequest("GET", "/product/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.getProduct(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("getProduct_Failed_notFound", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}
		c.Request, _ = http.NewRequest("GET", "/product/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.getProduct(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_productHandler_deleteProduct(t *testing.T) {
	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	id, catId := uuid.New(), uuid.New()
	mockProductRepo := &mockProductRepository{
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
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}
	t.Run("deleteProduct_Successful", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest("DELETE", "/product/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.deleteProduct(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, 0, len(mockProductRepo.items))
	})
	t.Run("deleteProduct_Failed_notValidId", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "id.String()"}}
		c.Request, _ = http.NewRequest("DELETE", "/product/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.deleteProduct(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, 0, len(mockProductRepo.items))
	})
	t.Run("deleteProduct_Failed_notFound", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}
		c.Request, _ = http.NewRequest("DELETE", "/product/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.deleteProduct(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, 0, len(mockProductRepo.items))
	})
}

func Test_productHandler_updateProduct(t *testing.T) {
	name, description := "test", "test"
	categoryName, categoryDescription := "test category", "test category description"
	id, catId := uuid.New(), uuid.New()

	mockProductRepo := &mockProductRepository{
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
	productHandler := &productHandler{
		productRepo: mockProductRepo,
	}

	t.Run("updateProduct_Successful", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/product/:id", nil)
		c.Params = []gin.Param{{Key: "id", Value: id.String()}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPUTPayload()))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.updateProduct(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "product name update", mockProductRepo.items[0].Name)
	})

	t.Run("updateProduct_Failed_idNotValid", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/product/:id", nil)
		c.Params = []gin.Param{{Key: "id", Value: "id.String()"}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPUTPayload()))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.updateProduct(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("updateProduct_Failed_invalidReqBody", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/product/:id", nil)
		c.Params = []gin.Param{{Key: "id", Value: id.String()}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid payload")))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.updateProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("updateProduct_Failed_invalidFormat", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/product/:id", nil)
		c.Params = []gin.Param{{Key: "id", Value: id.String()}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPOSTPayloadNotValid()))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.updateProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("updateProduct_Failed_invalidFormat", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("PUT", "/product/:id", nil)
		c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getProductPUTPayload()))
		c.Request.Header.Set("Content-Type", "application/json")
		productHandler.updateProduct(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

type mockProductRepository struct {
	items      []model.Product
	categories []model.Category
}

func (r *mockProductRepository) Insert(product *model.Product) error {
	newCategories := []model.Category{}
	for _, item := range product.Categories {
		for _, category := range r.categories {
			if item.ID == category.ID {
				newCategories = append(newCategories, category)
			}
		}
	}

	if len(newCategories) != len(product.Categories) {
		return fmt.Errorf("23503")
	}

	product.Categories = newCategories
	r.items = append(r.items, *product)

	return nil
}

// GetProducts get all products
func (r *mockProductRepository) GetAll(pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	var products []model.Product
	pagination.TotalRows = int64(len(r.items))
	pagination.Rows = ProductsToResponse(&products)

	return pagination, nil
}

// GetProduct get a single product
func (r *mockProductRepository) Get(id uuid.UUID) (*model.Product, error) {
	for _, item := range r.items {
		fmt.Println(item.ID)
		fmt.Println(id)
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// GetProductWithoutCategories get a single product
func (r *mockProductRepository) GetProductWithoutCategories(id uuid.UUID) (*model.Product, error) {
	return nil, nil
}

// DeleteProduct delete a single product
func (r *mockProductRepository) Delete(product *model.Product) error {
	for i, item := range r.items {
		if item.ID == product.ID {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

// UpdateProduct update a single product
func (r *mockProductRepository) Update(product *model.Product) error {
	for i, item := range r.items {
		if item.ID == product.ID {
			r.items[i] = *product
			return nil
		}
	}
	return gorm.ErrRecordNotFound
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
			return nil
		}
	}
	return gorm.ErrRecordNotFound
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
