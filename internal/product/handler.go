package product

import (
	"errors"
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"
)

type productHandler struct {
	productRepo *ProductRepository
}

func NewProductHandler(r *gin.RouterGroup, productRepo *ProductRepository, cfg *config.Config) {
	handler := &productHandler{productRepo: productRepo}

	r.GET("", handler.getProducts)
	r.GET("/:id", handler.getProduct)

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.POST("", handler.createProduct)
	r.PUT("/:id", handler.updateProduct)
	r.DELETE("/:id", handler.deleteProduct)
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
	c.JSON(201, ProductToResponse(product))
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

// getProduct gets a single product
func (r *productHandler) getProduct(c *gin.Context) {
	product := &model.Product{}

	if err := c.ShouldBindUri(&product); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	id := c.Param("id")

	product, err := r.productRepo.GetProduct(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"msg": "product not found"})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, ProductToResponse(product))
}

// deleteProduct deletes a single product
func (r *productHandler) deleteProduct(c *gin.Context) {
	product := &model.Product{}

	if err := c.ShouldBindUri(&product); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	id := c.Param("id")

	if err := r.productRepo.DeleteProduct(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"msg": "product not found"})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(204, nil)
}

// updateProduct updates a single product
func (r *productHandler) updateProduct(c *gin.Context) {
	productReq := &api.ProductRequest{}

	if err := c.ShouldBindUri(&productReq); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	product := ProductRequestToProduct(productReq)
	fmt.Println("************product: ", product)
	product.ID = strfmt.UUID(c.Param("id"))

	err := r.productRepo.UpdateProduct(product)
	// fmt.Println("-----------product: ", product.ToString())

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"msg": "product not found"})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("******************")
	fmt.Println("product: ", product.ToString())
	c.JSON(200, ProductToResponse(product))
}
