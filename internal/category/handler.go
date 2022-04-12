package category

import (
	"fmt"
	"patika-ecommerce/internal/api"
	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"
	file_helper "patika-ecommerce/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type categoryHandler struct {
	categoryRepo    *CategoryRepository
	categoryService *CategoryService
}

func NewCategoryHandler(r *gin.RouterGroup, cfg *config.Config, categoryService *CategoryService) {
	handler := &categoryHandler{
		categoryService: categoryService,
	}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.POST("", handler.createCategory)
	r.GET("", handler.getCategories)
	r.GET("/:id", handler.getCategory)
	r.PUT("/:id", handler.updateCategory)
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
	// TODO add pagination
	categories, err := r.categoryService.GetCategories()
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CategoriesToCategoryResponse(categories))
}

// getCategory returns a category
func (r *categoryHandler) getCategory(c *gin.Context) {
	category := &model.Category{}

	if err := c.ShouldBindUri(category); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	category, err := r.categoryService.GetCategoryByID(strfmt.UUID4(c.Param("id")))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CategoryToCategoryResponse(category))
}

// updateCategory updates a category
func (r *categoryHandler) updateCategory(c *gin.Context) {
	category := &model.Category{}

	if err := c.ShouldBindUri(category); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	reqBody := &api.CategoryRequest{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	category = CategoryRequestToCategory(reqBody)
	category.ID = strfmt.UUID(c.Param("id"))

	if err := r.categoryService.UpdateCategory(category); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CategoryToCategoryResponse(category))
}

// createBulkCategories creates a new categories with file upload
func (r *categoryHandler) createBulkCategories(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(400, gin.H{
			"message": "No file is received",
		})
		return
	}

	if err := file_helper.CheckFileIsValid(file); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	categories, err := r.categoryService.CreateBulkCategories(file)

	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Println(categories)
	fmt.Println(len(categories))
	c.JSON(200, CategoriesToCategoryResponse(&categories))
}
