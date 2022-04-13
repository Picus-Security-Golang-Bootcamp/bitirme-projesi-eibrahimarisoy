package product

import (
	"errors"
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/category"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"
	"patika-ecommerce/pkg/utils"

	mw "patika-ecommerce/pkg/middleware"

	httpErr "patika-ecommerce/internal/httpErrors"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"
)

type productHandler struct {
	productRepo  *ProductRepository
	categoryRepo *category.CategoryRepository
}

func NewProductHandler(r *gin.RouterGroup, cfg *config.Config, productRepo *ProductRepository, categoryRepo *category.CategoryRepository) {
	handler := &productHandler{productRepo: productRepo, categoryRepo: categoryRepo}
	// Public endpoints
	r.GET("", handler.getProducts)
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

	if err := c.ShouldBindUri(product); err != nil {
		c.JSON(httpErr.ErrorResponse(err)) // TODO payload error basiyor kontrol
		return
	}

	product, err := r.productRepo.GetProduct(strfmt.UUID(c.Param("id")))

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, ProductToResponse(product))
}

// deleteProduct deletes a single product
func (r *productHandler) deleteProduct(c *gin.Context) {
	product := &model.Product{}

	if err := c.ShouldBindUri(product); err != nil {
		c.JSON(httpErr.ErrorResponse(err)) // TODO payload error basiyor kontrol
		return
	}

	product.ID = strfmt.UUID(c.Param("id"))

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
	product.ID = strfmt.UUID(c.Param("id"))

	for index, item := range product.Categories {
		category, err := r.categoryRepo.GetCategoryByID(strfmt.UUID4(item.ID))

		if err != nil {
			c.JSON(httpErr.ErrorResponse(err))
			return
		}
		product.Categories[index] = *category
	}

	err := r.productRepo.UpdateProduct(product)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, ProductToResponse(product))
}

// FIXME
// ida := strfmt.UUID(c.Param("id"))

// 	id := c.Param("id")
// 	idx, a := uuid.Parse(ida.String())
// 	if !a {
// 		c.JSON(httpErr.ErrorResponse(errors.New("invalid product id")))
// 		return
// 	}

// [{{2022-04-10 11:21:50.094983 +0300 +03 2022-04-10 11:21:50.094983 +0300 +03 <nil>
// 	5183ee6e-6df5-429d-a897-7d4b77ffbdd8} 0xc0001600f0 name3 description3 []}
// 	{{2022-04-10 11:21:50.096588 +0300 +03 2022-04-10 11:21:50.096588 +0300 +03 <nil>
// 		267a5cad-35b3-46fa-9e03-619402ab7902} 0xc0001602b0 name4 description4 []}
// 		{{2022-04-10 11:21:50.093569 +0300 +03 2022-04-10 11:21:50.093569 +0300 +03 <nil>
// 			87b4e774-e879-4db9-8fd3-0c8bffcd1541} 0xc0001604b0 name2 description2 []}]
