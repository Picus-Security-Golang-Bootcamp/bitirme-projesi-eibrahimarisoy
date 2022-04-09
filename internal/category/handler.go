package category

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"
	file_helper "patika-ecommerce/pkg/utils"

	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryRepo    *CategoryRepository
	categoryService *CategoryService
}

func NewCategoryHandler(r *gin.RouterGroup, categoryRepo *CategoryRepository, categoryService *CategoryService, cfg *config.Config) {
	handler := &categoryHandler{
		categoryRepo:    categoryRepo,
		categoryService: categoryService,
	}
	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.POST("", handler.createCategory)
	r.GET("/", handler.getCategories)
	r.POST("/bulk-upload", handler.createBulkCategories)

	// r.GET("/:id", handler.getRole)
	// r.PUT("/:id", handler.updateRole)
	// r.DELETE("/:id", handler.deleteRole)
}

// createRole creates a new role
func (r *categoryHandler) createCategory(c *gin.Context) {
	categoryRequest := &api.CategoryRequest{}
	if err := c.ShouldBindJSON(&categoryRequest); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	category := CategoryRequestToCategory(categoryRequest)
	if err := r.categoryRepo.InsertCategory(category); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, CategoryToCategoryResponse(category))
}

// getCategory returns all roles
func (r *categoryHandler) getCategories(c *gin.Context) {

	categories, err := r.categoryRepo.GetCategories()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, CategoriesToCategoryResponse(categories))
}

// createBulkCategories
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
	c.JSON(200, gin.H{"message": "success"})
}

// _, exists := c.Get("user")
// if !exists {
// 	c.JSON(401, gin.H{"error": "Authentication required"})
// 	return
// }
