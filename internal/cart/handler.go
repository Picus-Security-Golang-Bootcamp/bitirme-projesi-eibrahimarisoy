package cart

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"

	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type cartHandler struct {
	cartService *CartService
}

func NewCartHandler(r *gin.RouterGroup, cfg *config.Config, cartService *CartService) {
	handler := &cartHandler{cartService: cartService}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("", handler.getOrCreateCart)
	r.POST("/add", handler.AddToCart)
	r.GET("/items", handler.ListCartItems)
	r.PUT("/items/:id", handler.UpdateCartItem)
	r.DELETE("/items/:id", handler.DeleteCartItem)
}

// createCart creates a new cart
func (r *cartHandler) getOrCreateCart(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.GetOrCreateCart(*user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, CartToCartResponse(cart))
}

// AddToCart adds a product to cart
func (r *cartHandler) AddToCart(c *gin.Context) {
	req := &api.CartAddRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.AddToCart(*user, req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, CartToCartResponse(cart))
}

// ListCartItem lists all cart items
func (r *cartHandler) ListCartItems(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.GetOrCreateCart(*user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, CartItemsToCartItemResponse(cart.Items))
}

// UpdateCartItem updates a cart item
func (r *cartHandler) UpdateCartItem(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	// TODO eger quantity 0 gelrise cikarma islemi yapilacak
	id := c.Param("id")
	req := &api.CartItemUpdateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cartItem, err := r.cartService.UpdateCartItem(*user, strfmt.UUID(id), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, CartItemToCartItemResponse(cartItem))
}

// DeleteCartItem deletes a cart item
func (r *cartHandler) DeleteCartItem(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	id := c.Param("id")

	err := r.cartService.DeleteCartItem(*user, strfmt.UUID(id))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Cart item deleted"})
}
