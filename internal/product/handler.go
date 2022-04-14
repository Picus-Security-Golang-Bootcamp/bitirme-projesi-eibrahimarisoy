package product

import (
	"errors"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/category"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"
	paginationHelper "patika-ecommerce/pkg/pagination"

	httpErr "patika-ecommerce/internal/httpErrors"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productHandler struct {
	productRepo  *ProductRepository
	categoryRepo *category.CategoryRepository
}

func NewProductHandler(r *gin.RouterGroup, cfg *config.Config, productRepo *ProductRepository, categoryRepo *category.CategoryRepository) {
	handler := &productHandler{productRepo: productRepo, categoryRepo: categoryRepo}
	// Public endpoints
	r.GET("", mw.PaginationMiddleware(), handler.getProducts)
	r.GET("/:id", handler.getProduct)

	// Private endpoints
	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.Use(mw.AdminMiddleware())
	r.POST("", handler.createProduct)
	r.PUT("/:id", handler.updateProduct)
	r.DELETE("/:id", handler.deleteProduct)
}

// createProduct creates a new product
func (r *productHandler) createProduct(c *gin.Context) {
	reqBody := &api.ProductRequest{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	product := ProductRequestToProduct(reqBody)

	if err := r.productRepo.InsertProduct(product); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(201, ProductToResponse(product))
}

// getProducts gets all products
func (r *productHandler) getProducts(c *gin.Context) {
	pagination := c.MustGet("pagination").(*paginationHelper.Pagination)

	data, err := r.productRepo.GetProducts(pagination)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, data)

}

// getProduct gets a single product
func (r *productHandler) getProduct(c *gin.Context) {
	product := &model.Product{}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	product, err = r.productRepo.GetProduct(id)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, ProductToResponse(product))
}

// deleteProduct deletes a single product
func (r *productHandler) deleteProduct(c *gin.Context) {
	product := &model.Product{}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	product.ID = id

	if err := r.productRepo.DeleteProduct(product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(httpErr.ErrorResponse(err))
			return
		}
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(204, nil)
}

// updateProduct updates a single product
func (r *productHandler) updateProduct(c *gin.Context) {
	reqBody := &api.ProductRequest{}

	if err := c.ShouldBindJSON(reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	product := ProductRequestToProduct(reqBody)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	product.ID = id

	for index, item := range product.Categories {
		category, err := r.categoryRepo.GetCategoryByID(item.ID)

		if err != nil {
			c.JSON(httpErr.ErrorResponse(err))
			return
		}
		product.Categories[index] = *category
	}

	if err = r.productRepo.UpdateProduct(product); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, ProductToResponse(product))
}
