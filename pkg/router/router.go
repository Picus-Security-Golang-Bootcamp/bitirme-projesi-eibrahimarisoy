package router

import (
	auth "patika-ecommerce/internal/auth"
	user "patika-ecommerce/internal/user"

	"patika-ecommerce/pkg/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitializeRoutes
func InitializeRoutes(rootRouter *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	authGroup := rootRouter.Group("/")
	// roleGroup := rootRouter.Group("/roles")
	// userGroup := rootRouter.Group("/users")
	// authGroup := rootRouter.Group("/auth")

	// Role repository
	// roleRepo := role.NewRoleRepository(db)
	// roleRepo.Migration()
	// role.NewRoleHandler(roleGroup, roleRepo)

	// // User repository
	userRepo := user.NewUserRepository(db)
	userRepo.Migration()
	// user.NewUserHandler(userGroup, userRepo)

	auth.NewAuthHandler(authGroup, userRepo)
}
