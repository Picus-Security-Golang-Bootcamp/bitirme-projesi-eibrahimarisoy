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
	"gorm.io/gorm"
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

	t.Run("createCategory_Succesfull", func(t *testing.T) {

		mockService := &mockCategoryService{
			items: []model.Category{},
		}
		categoryHandler := &categoryHandler{
			categoryService: mockService,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/categories", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getCategoryPOSTPayload()))
		categoryHandler.createCategory(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, 1, len(mockService.items))
	})

	t.Run("createCategory_Failed_reqBody", func(t *testing.T) {

		mockService := &mockCategoryService{
			items: []model.Category{},
		}
		categoryHandler := &categoryHandler{
			categoryService: mockService,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/categories", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("failed req body")))
		categoryHandler.createCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("createCategory_Failed_notValidate", func(t *testing.T) {
		payload := []byte(`{"name_fail":"test"}`)
		mockService := &mockCategoryService{
			items: []model.Category{},
		}
		categoryHandler := &categoryHandler{
			categoryService: mockService,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/categories", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(payload)))
		categoryHandler.createCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("createCategory_Failed_duplicateName", func(t *testing.T) {
		payload := []byte(`{"name":"test"}`)
		name := "test"
		mockService := &mockCategoryService{
			items: []model.Category{
				{
					Base: model.Base{ID: uuid.New()},
					Name: &name,
				},
			},
		}
		categoryHandler := &categoryHandler{
			categoryService: mockService,
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/categories", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(payload)))
		categoryHandler.createCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

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
	id := uuid.New()
	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}
	t.Run("getCategory_Succesfull", func(t *testing.T) {
		categoryHandler := &categoryHandler{
			categoryService: mockService,
		}
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest("GET", "categories/"+id.String(), nil)
		c.Request.Header.Set("Content-Type", "application/json")
		categoryHandler.getCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("getCategory_Failed_UUIDFault", func(t *testing.T) {
		categoryHandler := &categoryHandler{
			categoryService: mockService,
		}
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: "uuid-fault"}}
		c.Request, _ = http.NewRequest("GET", "categories/uuid-fault", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		categoryHandler.getCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
	t.Run("getCategory_Failed_notFound", func(t *testing.T) {
		categoryHandler := &categoryHandler{
			categoryService: mockService,
		}
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: uuid.New().String()}}
		c.Request, _ = http.NewRequest("GET", "categories/uuid-fault", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		categoryHandler.getCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_categoryHandler_updateCategory(t *testing.T) {
	name, description := "test", "test"
	id := uuid.New()

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	t.Run("updateCategory_Succesfull", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest("GET", "/categories/"+id.String(), nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getCategoryPUTPayload()))
		categoryHandler.updateCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("updateCategory_Failed_UUIDFault", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: "uuid-fault"}}
		c.Request, _ = http.NewRequest("GET", "/categories/"+"uuid-fault", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getCategoryPUTPayload()))
		categoryHandler.updateCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	payloadFault := []byte(`{"name": "test", 	}`)
	t.Run("updateCategory_Failed_reqBody", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest("GET", "/categories/"+id.String(), nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(payloadFault))
		categoryHandler.updateCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	payloadNotValid := []byte(`{"name_fault": "test", "description": "test", "id": "uuid-fault"}`)
	t.Run("updateCategory_Failed_notValidate", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest("GET", "/categories/"+id.String(), nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(payloadNotValid))
		categoryHandler.updateCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	notFoundUUID := uuid.New()
	t.Run("updateCategory_Failed_notFound", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: notFoundUUID.String()}}
		c.Request, _ = http.NewRequest("GET", "/categories/"+notFoundUUID.String(), nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getCategoryPUTPayload()))
		categoryHandler.updateCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

func Test_categoryHandler_deleteCategory(t *testing.T) {
	name, description := "test", "test"
	id := uuid.New()

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}

	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	t.Run("deleteCategory_Succesfull", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: id.String()}}
		r.DELETE("/categories/:id", categoryHandler.deleteCategory)
		c.Request, _ = http.NewRequest("GET", "/categories/"+id.String(), nil)
		c.Request.Header.Set("Content-Type", "application/json")
		categoryHandler.deleteCategory(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("deleteCategory_Failed_UUIDFault", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: "id.String()"}}
		r.DELETE("/categories/:id", categoryHandler.deleteCategory)
		c.Request, _ = http.NewRequest("GET", "/categories/"+"id.String()", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		categoryHandler.deleteCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	newId := uuid.New().String()
	t.Run("deleteCategory_Failed_UUIDFault", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, r := gin.CreateTestContext(w)
		c.Params = []gin.Param{gin.Param{Key: "id", Value: newId}}
		r.DELETE("/categories/:id", categoryHandler.deleteCategory)
		c.Request, _ = http.NewRequest("GET", "/categories/"+newId, nil)
		c.Request.Header.Set("Content-Type", "application/json")
		categoryHandler.deleteCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_categoryHandler_createBulkCategories(t *testing.T) {
	id, name, description := uuid.New(), "test", "test"

	mockService := &mockCategoryService{
		items: []model.Category{
			{
				Base:        model.Base{ID: id},
				Name:        &name,
				Description: description,
			},
		},
	}
	categoryHandler := &categoryHandler{
		categoryService: mockService,
	}
	t.Run("createBulkCategories_Succesfull", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "file.csv")
		part.Write([]byte("category_name,category_description\ncategory name 1,category description 1\ncategory name 2,category description 2"))
		writer.Close()

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		gin.SetMode(gin.TestMode)
		r.POST("/categories/bulk-upload", categoryHandler.createBulkCategories)
		c.Request, _ = http.NewRequest("GET", "/categories/bulk-upload", body)
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())
		categoryHandler.createBulkCategories(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("createBulkCategories_Failed_extension", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "file.pdf")
		part.Write([]byte("category_name,category_description\ncategory name 1,category description 1\ncategory name 2,category description 2"))
		writer.Close()

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		gin.SetMode(gin.TestMode)
		r.POST("/categories/bulk-upload", categoryHandler.createBulkCategories)
		c.Request, _ = http.NewRequest("GET", "/categories/bulk-upload", body)
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())
		categoryHandler.createBulkCategories(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("createBulkCategories_Failed_duplicateName", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "file.csv")
		part.Write([]byte("category_name,category_description\ncategory name 1,category description 1\ncategory name 1,category description 2"))
		writer.Close()

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		gin.SetMode(gin.TestMode)
		r.POST("/categories/bulk-upload", categoryHandler.createBulkCategories)
		c.Request, _ = http.NewRequest("GET", "/categories/bulk-upload", body)
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())
		categoryHandler.createBulkCategories(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

type mockCategoryService struct {
	items []model.Category
}

// CreateCategory creates a new category
func (c *mockCategoryService) CreateCategory(category *model.Category) error {
	for _, item := range c.items {
		if *item.Name == *category.Name {
			return errors.New("23505")
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
	return nil, gorm.ErrRecordNotFound
}

// UpdateCategory updates a category
func (c *mockCategoryService) UpdateCategory(category *model.Category) error {
	for index, item := range c.items {
		if item.ID == category.ID {
			c.items[index] = *category
			return nil
		}
	}
	return gorm.ErrRecordNotFound
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
				return nil, fmt.Errorf("23505")
			}
		}
		categories = append(categories, category)
	}

	for _, category := range categories {
		c.items = append(c.items, category)
	}

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
	return gorm.ErrRecordNotFound
}
