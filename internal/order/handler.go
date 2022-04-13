package order

import (
	"patika-ecommerce/internal/api"
	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"

	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type orderHandler struct {
	orderService *OrderService
}

func NewOrderHandler(r *gin.RouterGroup, cfg *config.Config, orderService *OrderService) {
	handler := &orderHandler{orderService: orderService}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("", handler.completeOrder)
	r.GET("/", handler.listOrders)
	r.PUT("/:id", handler.cancelOrder)
	//
}

// completeOrder completes an order
func (r *orderHandler) completeOrder(c *gin.Context) {
	reqBody := &api.OrderRequest{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	user := c.MustGet("user").(*model.User)

	order, err := r.orderService.CompleteOrder(user, reqBody)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, OrderToOrderResponse(order))
}

// listOrders lists all orders
func (r *orderHandler) listOrders(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	orders, err := r.orderService.GetOrdersByUser(user)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, OrdersToOrderDetailedResponse(orders))
}

// cancelOrder cancels an order
func (r *orderHandler) cancelOrder(c *gin.Context) {
	id := c.Param("id")

	user := c.MustGet("user").(*model.User)

	err := r.orderService.CancelOrder(user, strfmt.UUID(id))

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(200, map[string]string{"message": "order cancelled"})
}
