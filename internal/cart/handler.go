package cart

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"

	"patika-ecommerce/pkg/config"
	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type cartHandler struct {
	cartService *CartService
}

func NewCartHandler(r *gin.RouterGroup, cfg *config.Config, cartService *CartService) {
	handler := &cartHandler{cartService: cartService}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("", handler.getOrCreateCart)
	r.POST("/add", handler.AddToCart)
	r.GET("/list", handler.ListCartItem)
	// r.GET("", handler.getCarts)
	// r.GET("/:id", handler.getCart)
	// r.PUT("/:id", handler.updateCart)
	// r.DELETE("/:id", handler.deleteCart)
}

// createCart creates a new cart
func (r *cartHandler) getOrCreateCart(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	fmt.Println(user)
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
func (r *cartHandler) ListCartItem(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	cart, err := r.cartService.GetOrCreateCart(*user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, CartItemsToCartItemResponse(cart.Items))
}
