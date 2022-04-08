package router

import (
	auth "patika-ecommerce/internal/auth"
	category "patika-ecommerce/internal/category"
	"patika-ecommerce/internal/role"
	user "patika-ecommerce/internal/user"

	"patika-ecommerce/pkg/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitializeRoutes
func InitializeRoutes(rootRouter *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	authGroup := rootRouter.Group("/")
	roleGroup := rootRouter.Group("/roles")
	userGroup := rootRouter.Group("/users")
	categoryGroup := rootRouter.Group("/categories")

	// Role repository
	roleRepo := role.NewRoleRepository(db)
	roleRepo.Migration()
	role.NewRoleHandler(roleGroup, roleRepo, cfg)

	// // User repository
	userRepo := user.NewUserRepository(db)
	userRepo.Migration()
	user.NewUserHandler(userGroup, userRepo)

	authService := auth.NewAuthService(cfg)
	auth.NewAuthHandler(authGroup, userRepo, authService)

	// Category repository
	categoryRepo := category.NewCategoryRepository(db)
	categoryRepo.Migration()
	category.NewCategoryHandler(categoryGroup, categoryRepo)
}
