package router

import (
	auth "patika-ecommerce/internal/auth"
	cart "patika-ecommerce/internal/cart"
	category "patika-ecommerce/internal/category"
	"patika-ecommerce/internal/order"
	product "patika-ecommerce/internal/product"
	user "patika-ecommerce/internal/user"

	"patika-ecommerce/pkg/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitializeRoutes
func InitializeRoutes(rootRouter *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {

	// Initialize the router groups
	authGroup := rootRouter.Group("/")
	categoryGroup := rootRouter.Group("/categories")
	productGroup := rootRouter.Group("/products")
	cartGroup := rootRouter.Group("/cart")
	orderGroup := rootRouter.Group("/orders")

	// User repository
	userRepo := user.NewUserRepository(db)
	userRepo.Migration()
	// Auth service
	authService := auth.NewAuthService(cfg, userRepo)
	auth.NewAuthHandler(authGroup, cfg, authService)

	// Category repository
	categoryRepo := category.NewCategoryrRepository(db)
	categoryRepo.Migration()
	categoryService := category.NewCategoryService(categoryRepo)
	category.NewCategoryHandler(categoryGroup, cfg, categoryService)

	// Product repository
	productRepo := product.NewProductRepository(db)
	productRepo.Migration()
	product.NewProductHandler(productGroup, cfg, productRepo)

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
	order.NewOrderHandler(orderGroup, cfg, orderRepo)

}
