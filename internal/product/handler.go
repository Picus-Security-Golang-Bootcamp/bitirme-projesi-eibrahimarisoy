package product

import (
	"errors"
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"
	"patika-ecommerce/pkg/utils"

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
	pagination := utils.GeneratePaginationFromRequest(c)
	data, totalPages, err := r.productRepo.GetProducts(&pagination)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// get current url path
	urlPath := c.Request.URL.Path

	// search query params
	searchQueryParams := pagination.Q

	// set first & last page pagination response
	data.FirstPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&q=%s", urlPath, pagination.Limit, 0, pagination.Sort, searchQueryParams)
	data.LastPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&q=%s", urlPath, pagination.Limit, totalPages, pagination.Sort, searchQueryParams)

	fmt.Println("data: ", data.Page)
	if data.Page > 1 {
		// set previous page pagination response
		data.PreviousPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&q=%s", urlPath, pagination.Limit, data.Page-1, pagination.Sort, searchQueryParams)
	}

	if data.Page < totalPages {
		// set next page pagination response
		data.NextPage = fmt.Sprintf("%s?limit=%d&page=%d&sort=%s&q=%s", urlPath, pagination.Limit, data.Page+1, pagination.Sort, searchQueryParams)
	}

	if data.Page > totalPages {
		// reset previous page
		data.PreviousPage = ""
	}

	// if data.Page == totalPages {
	// 	// reset next page
	// 	data.NextPage = ""
	// }

	c.JSON(200, data)
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
