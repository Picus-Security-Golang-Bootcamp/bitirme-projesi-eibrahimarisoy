package cart

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"

	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type cartHandler struct {
	cartService *CartService
}

// NewCartHandler creates a new cart handler
func NewCartHandler(r *gin.RouterGroup, cfg *config.Config, cartService *CartService) {
	handler := &cartHandler{cartService: cartService}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("", handler.getOrCreateCart)
	r.POST("/add", handler.AddToCart)
	r.GET("/items", handler.ListCartItems)
	r.PUT("/items/:id", handler.UpdateCartItem)
	r.DELETE("/items/:id", handler.DeleteCartItem)
}

// getOrCreateCart if users cart exists, returns it, otherwise creates a new one
func (r *cartHandler) getOrCreateCart(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.GetOrCreateCart(user)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CartToCartResponse(cart))
}

// AddToCart adds a product to cart
func (r *cartHandler) AddToCart(c *gin.Context) {
	reqBody := &api.AddToCartRequest{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.AddToCart(user, reqBody)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CartToCartResponse(cart))
}

// ListCartItem lists all cart items
func (r *cartHandler) ListCartItems(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.GetOrCreateCart(user)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CartItemsToCartItemResponse(cart.Items))
}

// UpdateCartItem updates a cart item
func (r *cartHandler) UpdateCartItem(c *gin.Context) {
	reqBody := &api.CartItemUpdateRequest{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	user := c.MustGet("user").(*model.User)

	cartItem, err := r.cartService.UpdateCartItem(user, strfmt.UUID(c.Param("id")), reqBody)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(200, CartItemToCartItemResponse(cartItem))
}

// DeleteCartItem deletes a cart item
func (r *cartHandler) DeleteCartItem(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	err := r.cartService.DeleteCartItem(user, strfmt.UUID(c.Param("id")))

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(204, nil)
}
