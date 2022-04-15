package category

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func getCategoryPOSTPayload() []byte {
	var jsonStr = []byte(
		`{"name":"test","description":"description"}`)

	return jsonStr
}

func getCategoryPUTPayload() []byte {
	var jsonStr = []byte(
		`{"name":"test update","description":"description update"}`)

	return jsonStr
}

func Test_categoryHandler_createCategory(t *testing.T) {
	// name, description := "test", "test"
	// id := uuid.New()

	gin.SetMode(gin.TestMode)

	mockService := &mockCategoryService{
		items: []model.Category{
			// {
			// 	Base:        model.Base{ID: id},
			// 	Name:        &name,
			// 	Description: description,
			// },
		},
	}
	w := httptest.NewRecorder()
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	c, r := gin.CreateTestContext(w)

	r.POST("/categories", categoryHandler.createCategory)
	c.Request, _ = http.NewRequest("POST", "/categories", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getCategoryPOSTPayload()))
	categoryHandler.createCategory(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	fmt.Println(len(mockService.items))
	assert.Equal(t, 1, len(mockService.items))

}

func Test_categoryHandler_getCategories(t *testing.T) {
	name, description := "test", "test"
	name2, description2 := "test2", "test2"
	id := uuid.New()
	id2 := uuid.New()
	gin.SetMode(gin.TestMode)

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
			{
				Base:        model.Base{ID: id2},
				Name:        &name2,
				Description: description2,
			},
		},
	}
	w := httptest.NewRecorder()
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	c, r := gin.CreateTestContext(w)

	r.GET("/categories", categoryHandler.getCategories)
	c.Request, _ = http.NewRequest("GET", "/categories", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	categoryHandler.getCategories(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(mockService.items))
	assert.Equal(t, name, mockService.items[0].Name)
}

func Test_categoryHandler_getCategory(t *testing.T) {
	name, description := "test", "test"
	idx := uuid.New()

	gin.SetMode(gin.TestMode)

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: idx},
				Name:        &name,
				Description: description,
			},
		},
	}
	w := httptest.NewRecorder()
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	c, r := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: idx.String()}}
	r.GET("/:id", categoryHandler.getCategory)

	url := fmt.Sprintf("/%s", idx.String())
	fmt.Println(url)
	c.Request, _ = http.NewRequest("GET", url, nil)
	c.Request.Header.Set("Content-Type", "application/json")
	categoryHandler.getCategory(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_categoryHandler_updateCategory(t *testing.T) {
	name, description := "test", "test"
	id := uuid.New()
	fmt.Println(id)

	gin.SetMode(gin.TestMode)

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}
	w := httptest.NewRecorder()
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	c, r := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: id.String()}}

	r.PUT("/categories/:id", categoryHandler.updateCategory)
	c.Request, _ = http.NewRequest("GET", "/categories/"+id.String(), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getCategoryPUTPayload()))
	categoryHandler.updateCategory(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_categoryHandler_deleteCategory(t *testing.T) {
	name, description := "test", "test"
	id := uuid.New()

	gin.SetMode(gin.TestMode)

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}
	w := httptest.NewRecorder()
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	c, r := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "id", Value: id.String()}}

	r.DELETE("/categories/:id", categoryHandler.deleteCategory)
	c.Request, _ = http.NewRequest("GET", "/categories/"+id.String(), nil)
	c.Request.Header.Set("Content-Type", "application/json")
	categoryHandler.deleteCategory(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func Test_categoryHandler_createBulkCategories(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// writer.WriteField("bu", "HFL")
	// writer.WriteField("wk", "10")
	part, _ := writer.CreateFormFile("file", "file.csv")
	part.Write([]byte("category_name,category_description\ncategory name 1,category description 1\ncategory name 2,category description 2"))
	writer.Close()

	name, description := "test", "test"
	id := uuid.New()

	gin.SetMode(gin.TestMode)

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}
	w := httptest.NewRecorder()
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	c, r := gin.CreateTestContext(w)

	r.POST("/categories/bulk-upload", categoryHandler.createBulkCategories)
	c.Request, _ = http.NewRequest("GET", "/categories/bulk-upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())
	categoryHandler.createBulkCategories(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

type mockCategoryService struct {
	items []model.Category
}

// CreateCategory creates a new category
func (c *mockCategoryService) CreateCategory(category *model.Category) error {
	for _, item := range c.items {
		if *item.Name == *category.Name {
			return errors.New("category already exists")
		}
	}
	c.items = append(c.items, *category)

	return nil
}

// GetCategories returns all categories
func (c *mockCategoryService) GetCategories() (*[]model.Category, error) {

	return &c.items, nil
}

// GetCategoryByID returns a category by id
func (c *mockCategoryService) GetCategoryByID(id uuid.UUID) (*model.Category, error) {

	for _, item := range c.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, nil
}

// UpdateCategory updates a category
func (c *mockCategoryService) UpdateCategory(category *model.Category) error {
	for index, item := range c.items {
		if item.ID == category.ID {
			c.items[index] = *category
			return nil
		}
	}
	return errors.New("category not found")
}

// CreateBulkCategories creates multiple categories in bulk operation with the specified file
func (c *mockCategoryService) CreateBulkCategories(filename *bytes.Buffer) ([]model.Category, error) {
	records, err := utils.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var categories []model.Category
	for _, record := range records[1:] {
		category := model.Category{
			Name:        &record[0],
			Description: record[1],
		}
		for _, item := range c.items {
			if *item.Name == *category.Name {
				return nil, fmt.Errorf("category already exists")
			}
		}
		categories = append(categories, category)
	}

	for _, category := range categories {
		c.items = append(c.items, category)
	}

	// if err := c.categoryRepo.InsertBulkCategory(&categories); err != nil {
	// 	return nil, err
	// }
	return categories, nil

}

// DeleteCategory deletes a category by id
func (c *mockCategoryService) DeleteCategoryService(id uuid.UUID) error {

	for index, item := range c.items {
		if item.ID == id {
			c.items = append(c.items[:index], c.items[index+1:]...)
			return nil
		}
	}
	return errors.New("category not found")
}
