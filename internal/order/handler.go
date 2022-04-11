package order

import (
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderService *OrderService
}

func NewOrderHandler(r *gin.RouterGroup, cfg *config.Config, orderService *OrderService) {
	// handler := &orderHandler{orderService: orderService}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	// r.POST("", handler.getOrCreateOrder)
	//
}
