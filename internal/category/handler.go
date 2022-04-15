package category

import (
	"bytes"
	"fmt"
	"io"
	"patika-ecommerce/internal/api"
	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"
	file_helper "patika-ecommerce/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type categoryHandler struct {
	categoryService MockCategoryService
}

func NewCategoryHandler(r *gin.RouterGroup, cfg *config.Config, categoryService MockCategoryService) {
	handler := &categoryHandler{
		categoryService: categoryService,
	}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.POST("", handler.createCategory)
	r.GET("", handler.getCategories)
	r.GET("/:id", handler.getCategory)
	r.PUT("/:id", handler.updateCategory)
	r.DELETE("/:id", handler.deleteCategory)
	r.POST("/bulk-upload", handler.createBulkCategories)
}

// createCategory creates a new category
func (r *categoryHandler) createCategory(c *gin.Context) {
	reqBody := &api.CategoryRequest{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	category := CategoryRequestToCategory(reqBody)

	if err := r.categoryService.CreateCategory(category); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(201, CategoryToCategoryResponse(category))
}

// getCategories returns all categories
func (r *categoryHandler) getCategories(c *gin.Context) {
	categories, err := r.categoryService.GetCategories()
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CategoriesToCategoryResponse(categories))
}

// getCategory returns a category
func (r *categoryHandler) getCategory(c *gin.Context) {
	fmt.Println("category update", c.Param("id"))

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	fmt.Println(id)

	category, err := r.categoryService.GetCategoryByID(id)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CategoryToCategoryResponse(category))
}

// updateCategory updates a category
func (r *categoryHandler) updateCategory(c *gin.Context) {
	fmt.Println("category update", c.Param("id"))
	categoryID, err := uuid.Parse(c.Param("id"))
	fmt.Println(categoryID)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	fmt.Println("category", c.Param("id"))
	reqBody := &api.CategoryRequest{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	category := CategoryRequestToCategory(reqBody)
	fmt.Println(category)
	category.ID = categoryID

	if err := r.categoryService.UpdateCategory(category); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CategoryToCategoryResponse(category))
}

// createBulkCategories creates a new categories with file upload
func (r *categoryHandler) createBulkCategories(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	// defer file.Close()

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
	}

	if err := file_helper.CheckFileIsValid(header); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
	}

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	categories, err := r.categoryService.CreateBulkCategories(buf)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CategoriesToCategoryResponse(&categories))
}

// deleteCategory deletes a category
func (r *categoryHandler) deleteCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := r.categoryService.DeleteCategoryService(categoryID); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(204, nil)
}
