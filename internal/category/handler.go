package category

import (
	"encoding/csv"
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"
	file_helper "patika-ecommerce/pkg/utils"

	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	categoryRepo *CategoryRepository
}

func NewCategoryHandler(r *gin.RouterGroup, categoryRepo *CategoryRepository, cfg *config.Config) {
	handler := &categoryHandler{categoryRepo: categoryRepo}
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

	f, _ := file.Open()
	defer f.Close()

	lines, _ := csv.NewReader(f).ReadAll()
	for _, line := range lines[1:] {
		fmt.Println(line)
	}

	c.JSON(200, gin.H{"message": "success"})
}

// _, exists := c.Get("user")
// if !exists {
// 	c.JSON(401, gin.H{"error": "Authentication required"})
// 	return
// }
