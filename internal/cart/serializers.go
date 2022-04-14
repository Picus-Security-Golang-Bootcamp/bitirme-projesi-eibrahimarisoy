package cart

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/internal/product"
	common "patika-ecommerce/pkg/utils"
)

func CartToCartResponse(cart *model.Cart) *api.CartResponse {
	cartItemResponse := []*api.CartItemResponse{}

	for _, v := range cart.Items {
		item := &api.CartItemResponse{
			ID:       common.UUIDToStrfmt(v.ID),
			Product:  common.UUIDToStrfmt(v.ProductID),
			Quantity: int64(v.Quantity),
			Price:    v.Price,
		}
		cartItemResponse = append(cartItemResponse, item)
	}

	return &api.CartResponse{
		ID:         common.UUIDToStrfmt(cart.ID),
		Status:     string(cart.Status),
		Items:      cartItemResponse,
		TotalPrice: cart.GetTotalPrice(),
	}
}

func CartAddRequestToCartItem(req *api.AddToCartRequest) *model.CartItem {
	id, _ := common.StrfmtToUUID(req.ProductID)

	return &model.CartItem{
		ProductID: id,
		Quantity:  req.Quantity,
	}
}

// CartItemToCartItemResponse converts a cart item to a cart item response
func CartItemToCartItemResponse(item *model.CartItem) *api.CartItemDetailResponse {

	return &api.CartItemDetailResponse{
		ID:       common.UUIDToStrfmt(item.ID),
		Product:  product.ProductToProductBasicResponse(&item.Product),
		Quantity: int64(item.Quantity),
		Price:    item.Price,
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
