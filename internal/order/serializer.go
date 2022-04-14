package order

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/internal/product"
	common "patika-ecommerce/pkg/utils"
	"strconv"

	"github.com/go-openapi/strfmt"
)

// OrderToOrderResponse converts an order to an order response
func OrderToOrderResponse(order *model.Order) *api.OrderResponse {
	totalPrice, _ := strconv.ParseFloat(order.TotalPrice.String(), 64)

	return &api.OrderResponse{
		ID:         common.UUIDToStrfmt(order.ID),
		CartID:     common.UUIDToStrfmt(order.CartID),
		Status:     string(order.Status),
		TotalPrice: totalPrice,
		CreatedAt:  strfmt.DateTime(order.CreatedAt),
		UpdatedAt:  strfmt.DateTime(order.UpdatedAt),
	}
}

// OrderToOrderDetailedResponse converts an order to an order response
func OrderToOrderDetailedResponse(order *model.Order) *api.OrderDetailedResponse {
	items := []*api.OrderItemDetailedResponse{}

	for _, order := range order.Items {
		items = append(items, OrderItemToOrderItemDetailedResponse(&order))
	}
	totalPrice, _ := strconv.ParseFloat(order.TotalPrice.String(), 64)

	return &api.OrderDetailedResponse{
		ID:         common.UUIDToStrfmt(order.ID),
		CartID:     common.UUIDToStrfmt(order.CartID),
		Status:     string(order.Status),
		TotalPrice: totalPrice,
		Items:      items,
		CreatedAt:  strfmt.DateTime(order.CreatedAt),
		UpdatedAt:  strfmt.DateTime(order.UpdatedAt),
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
	price, _ := strconv.ParseFloat(orderItem.Price.String(), 64)

	return &api.OrderItemDetailedResponse{
		ID:      common.UUIDToStrfmt(orderItem.ID),
		Product: product.ProductToProductBasicResponse(&orderItem.Product),
		Price:   price,
	}
}
