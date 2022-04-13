package order

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/internal/product"
)

// OrderToOrderResponse converts an order to an order response
func OrderToOrderResponse(order *model.Order) *api.OrderResponse {
	return &api.OrderResponse{
		ID:         order.ID,
		CartID:     order.CartID,
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
		// CreatedAt: order.CreatedAt, // TODO add created at
		// UpdatedAt: order.UpdatedAt,
	}
}

// OrderToOrderDetailedResponse converts an order to an order response
func OrderToOrderDetailedResponse(order *model.Order) *api.OrderDetailedResponse {
	items := []*api.OrderItemDetailedResponse{}

	for _, order := range order.Items {
		items = append(items, OrderItemToOrderItemDetailedResponse(&order))
	}

	return &api.OrderDetailedResponse{
		ID:         order.ID,
		CartID:     order.CartID,
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
		Items:      items,
		// CreatedAt: order.CreatedAt, // TODO add created at
	}
}

// OrdersToOrderDetaledResponse converts an orders to an orders response
func OrdersToOrderDetailedResponse(orders []*model.Order) []*api.OrderDetailedResponse {
	var orderResponses []*api.OrderDetailedResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, OrderToOrderDetailedResponse(order))
	}
	return orderResponses
}

// OrderItemToOrderItemDetailedResponse converts an order item to an order item response
func OrderItemToOrderItemDetailedResponse(orderItem *model.OrderItem) *api.OrderItemDetailedResponse {
	return &api.OrderItemDetailedResponse{
		ID: orderItem.ID,
		// Quantity: orderItem.Quantity, // TODO remove quantity
		Product: product.ProductToProductBasicResponse(&orderItem.Product),
		Price:   orderItem.Price,
	}
}
