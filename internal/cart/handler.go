package cart

import (
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
	r.POST("", handler.createCart)

	// r.GET("", handler.getCarts)
	// r.GET("/:id", handler.getCart)
	// r.PUT("/:id", handler.updateCart)
	// r.DELETE("/:id", handler.deleteCart)
}

// createCart creates a new cart
func (r *cartHandler) createCart(c *gin.Context) {
	// cartReq := api.CartRequest{}
	// if err := c.ShouldBindJSON(&cartReq); err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }
	// fmt.Println("cart: ", &cartReq)
	// cart := CartRequestToCart(&cartReq)
	// if err := r.cartService.InsertCart(cart); err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }
	// c.JSON(201, CartToResponse(cart))
	c.JSON(201, gin.H{"message": "created cart"})
}
