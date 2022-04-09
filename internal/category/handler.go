package category

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"

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

	// userId := user.(*jwt_helper.JWTToken).UserId
	// fmt.Println(userId)

	// categories, err := r.categoryRepo.GetCategories()
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(200, nil)
}

// _, exists := c.Get("user")
// if !exists {
// 	c.JSON(401, gin.H{"error": "Authentication required"})
// 	return
// }
