package cart

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/internal/product"
)

func CartToCartResponse(cart *model.Cart) *api.CartResponse {
	cartItemResponse := []*api.CartItemResponse{}

	for _, v := range cart.Items {

		item := &api.CartItemResponse{
			ID:       v.ID,
			Product:  v.ProductID,
			Quantity: int64(v.Quantity),
			Price:    float64(v.Price),
		}
		cartItemResponse = append(cartItemResponse, item)
	}

	return &api.CartResponse{
		ID:     cart.ID,
		Status: string(cart.Status),
		Items:  cartItemResponse,
	}
}

func CartAddRequestToCartItem(req *api.CartAddRequest) *model.CartItem {
	return &model.CartItem{
		ProductID: req.ProductID,
		Quantity:  int(req.Quantity),
	}
}

// CartItemToCartItemResponse converts a cart item to a cart item response
func CartItemToCartItemResponse(item *model.CartItem) *api.CartItemDetailResponse {
	return &api.CartItemDetailResponse{
		ID:       item.ID,
		Product:  product.ProductToProductBasicResponse(&item.Product),
		Quantity: int64(item.Quantity),
		Price:    float64(item.Price),
	}
}

// CartItemsToCartItemResponse converts a cart item to a cart item response
func CartItemsToCartItemResponse(items []model.CartItem) []*api.CartItemDetailResponse {
	cartItemResponse := []*api.CartItemDetailResponse{}

	for _, v := range items {
		item := CartItemToCartItemResponse(&v)
		cartItemResponse = append(cartItemResponse, item)
	}

	return cartItemResponse
}
