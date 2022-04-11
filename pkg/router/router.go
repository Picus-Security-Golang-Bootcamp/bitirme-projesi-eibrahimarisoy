package router

import (
	auth "patika-ecommerce/internal/auth"
	cart "patika-ecommerce/internal/cart"
	category "patika-ecommerce/internal/category"
	"patika-ecommerce/internal/order"
	product "patika-ecommerce/internal/product"
	role "patika-ecommerce/internal/role"
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
	productGroup := rootRouter.Group("/products")
	cartGroup := rootRouter.Group("/cart")
	orderGroup := rootRouter.Group("/orders")

	// Role repository
	roleRepo := role.NewRoleRepository(db)
	roleRepo.Migration()
	role.NewRoleHandler(roleGroup, roleRepo, cfg)

	// // User repository
	userRepo := user.NewUserRepository(db)
	userRepo.Migration()
	user.NewUserHandler(userGroup, userRepo)

	// Auth service
	authService := auth.NewAuthService(cfg)
	auth.NewAuthHandler(authGroup, userRepo, authService)

	// Category repository
	categoryRepo := category.NewCategoryrRepository(db)
	categoryRepo.Migration()
	categoryService := category.NewCategoryService(categoryRepo)
	category.NewCategoryHandler(categoryGroup, categoryRepo, categoryService, cfg)

	// Product repository
	productRepo := product.NewProductRepository(db)
	productRepo.Migration()
	product.NewProductHandler(productGroup, productRepo, cfg)

	// Cart repository
	cartRepo := cart.NewCartRepository(db)
	cartRepo.Migration()
	cartItemRepo := cart.NewCartItemRepository(db)
	cartItemRepo.Migration()
	cartService := cart.NewCartService(cartRepo, productRepo, cartItemRepo)
	cart.NewCartHandler(cartGroup, cfg, cartService)

	// Order repository
	orderRepo := order.NewOrderRepository(db)
	orderRepo.Migration()
	orderItemRepo := order.NewOrderItemRepository(db)
	orderItemRepo.Migration()
	orderService := order.NewOrderService(orderRepo, orderItemRepo, cartRepo)
	order.NewOrderHandler(orderGroup, cfg, orderService)

}
