package category

import (
	"patika-ecommerce/pkg/config"

	"github.com/gin-gonic/gin"
)

type roleHandler struct {
	roleRepo *CategoryRepository
}

func NewRoleHandler(r *gin.RouterGroup, cfg *config.Config) {
	// handler := &roleHandler{roleRepo: roleRepo}
	// r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	// r.Use(mw.AdminMiddleware())
	// r.POST("", handler.createRole)
	// r.GET("/", handler.getRoles)
	// r.GET("/:id", handler.getRole)
	// r.PUT("/:id", handler.updateRole)
	// r.DELETE("/:id", handler.deleteRole)
}

// createRole creates a new role
// func (r *roleHandler) createRole(c *gin.Context) {
// 	role := &model.Role{}
// 	if err := c.ShouldBindJSON(&role); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := r.roleRepo.InsertRole(role); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(200, role)
// }
