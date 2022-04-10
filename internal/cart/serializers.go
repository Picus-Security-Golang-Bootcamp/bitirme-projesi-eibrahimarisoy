package cart

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
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
