package cart

import (
	"fmt"
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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, CartToCartResponse(cart))
}
