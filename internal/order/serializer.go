package order

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
)

// OrderToOrderResponse converts an order to an order response
func OrderToOrderResponse(order *model.Order) *api.OrderResponse {
	return &api.OrderResponse{
		ID:     order.ID,
		CartID: order.CartID,
		Status: string(order.Status),
		// CreatedAt: order.CreatedAt, // TODO add created at
		// UpdatedAt: order.UpdatedAt,
	}
}
