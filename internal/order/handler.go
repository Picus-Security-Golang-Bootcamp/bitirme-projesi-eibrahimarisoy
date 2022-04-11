package order

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderService *OrderService
}

func NewOrderHandler(r *gin.RouterGroup, cfg *config.Config, orderService *OrderService) {
	handler := &orderHandler{orderService: orderService}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("", handler.completeOrder)
	//
}

// completeOrder completes an order
func (r *orderHandler) completeOrder(c *gin.Context) {
	req := &api.OrderRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*model.User)

	order, err := r.orderService.CompleteOrder(user, req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, order)
}
