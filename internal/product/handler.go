package product

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	productRepo *ProductRepository
}

func NewProductHandler(r *gin.RouterGroup, productRepo *ProductRepository, cfg *config.Config) {
	handler := &productHandler{productRepo: productRepo}
	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.POST("", handler.createProduct)
	r.GET("", handler.getProducts)
	// r.GET("/:id", handler.getProduct)
	// r.PUT("/:id", handler.updateProduct)
	// r.DELETE("/:id", handler.deleteroduct)
}

// createProduct creates a new product
func (r *productHandler) createProduct(c *gin.Context) {
	productReq := api.ProductRequest{}
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("product: ", &productReq)
	product := ProductRequestToProduct(&productReq)
	if err := r.productRepo.InsertProduct(product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ProductToResponse(product))
}

// getProducts gets all products
func (r *productHandler) getProducts(c *gin.Context) {
	products, err := r.productRepo.GetProducts()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ProductsToResponse(products))
}
