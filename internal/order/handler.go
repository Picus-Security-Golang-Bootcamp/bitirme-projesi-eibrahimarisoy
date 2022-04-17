package order

import (
	"patika-ecommerce/internal/api"
	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"
	paginationHelper "patika-ecommerce/pkg/pagination"
	common "patika-ecommerce/pkg/utils"

	mw "patika-ecommerce/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

type orderHandler struct {
	orderRepo MockOrderRepository
}

func NewOrderHandler(r *gin.RouterGroup, cfg *config.Config, orderRepo *OrderRepository) {
	handler := &orderHandler{orderRepo: orderRepo}

	r.Use(mw.AuthenticationMiddleware(cfg.JWTConfig.SecretKey))
	r.POST("", handler.completeOrder)
	r.GET("", mw.PaginationMiddleware(), handler.listOrders)
	r.PUT("/:id", handler.cancelOrder)
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

	cartId, err := common.StrfmtToUUID(*reqBody.CartID)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
	}

	order, err := r.orderRepo.CompleteOrder(user, cartId)
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, OrderToOrderResponse(order))
}

// listOrders lists all orders
func (r *orderHandler) listOrders(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	pagination := c.MustGet("pagination").(*paginationHelper.Pagination)

	data, err := r.orderRepo.GetOrdersByUser(user, pagination)

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, data)
}

// cancelOrder cancels an order
func (r *orderHandler) cancelOrder(c *gin.Context) {
	user := c.MustGet("user").(*model.User)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := r.orderRepo.CancelOrder(id, user); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	c.JSON(200, map[string]string{"message": "order cancelled"})
}
