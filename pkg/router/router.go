package router

import (
	"patika-ecommerce/internal/role"
	"patika-ecommerce/pkg/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitializeRoutes
func InitializeRoutes(rootRouter *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {

	roleGroup := rootRouter.Group("/roles")
	// userGroup := rootRouter.Group("/users")
	// authGroup := rootRouter.Group("/auth")

	// Role repository
	roleRepo := role.NewRoleRepository(db)
	roleRepo.Migration()
	role.NewRoleHandler(roleGroup, roleRepo)

}
