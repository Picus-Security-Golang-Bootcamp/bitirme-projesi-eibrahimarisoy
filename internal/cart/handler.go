package cart

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"

	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type cartHandler struct {
	cartService MockCartService
}

// NewCartHandler creates a new cart handler
func NewCartHandler(r *gin.RouterGroup, cfg *config.Config, cartService MockCartService) {
	handler := &cartHandler{cartService: cartService}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("", handler.getOrCreateCart)
	r.POST("/add", handler.addToCart)
	r.GET("/items", handler.listCartItems)
	r.PUT("/items/:id", handler.updateCartItem)
	r.DELETE("/items/:id", handler.deleteCartItem)
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
func (r *cartHandler) addToCart(c *gin.Context) {
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
func (r *cartHandler) listCartItems(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.GetOrCreateCart(user)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, CartItemsToCartItemResponse(cart.Items))
}

// UpdateCartItem updates a cart item
func (r *cartHandler) updateCartItem(c *gin.Context) {
	reqBody := &api.CartItemUpdateRequest{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		fmt.Println("err", err)
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	user := c.MustGet("user").(*model.User)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	cartItem, err := r.cartService.UpdateCartItem(user, id, reqBody)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(200, CartItemToCartItemResponse(cartItem))
}

// DeleteCartItem deletes a cart item
func (r *cartHandler) deleteCartItem(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err = r.cartService.DeleteCartItem(user, id); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(204, nil)
}
