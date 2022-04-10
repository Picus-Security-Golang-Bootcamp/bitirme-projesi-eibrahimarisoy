package cart

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	product "patika-ecommerce/internal/product"
)

func CartToCartResponse(cart *model.Cart) *api.CartResponse {
	cartItemResponse := []*api.CartItemResponse{}

	for _, v := range cart.Items {

		item := &api.CartItemResponse{
			ID:       v.ID,
			Product:  product.ProductToResponse(&v.Product),
			Quantity: int64(v.Quantity),
			Price:    float64(v.Price),
		}

		cartItemResponse = append(cartItemResponse, item)
	}
	fmt.Println(cart.Status)
	fmt.Println(string(cart.Status))

	return &api.CartResponse{
		ID:     cart.ID,
		Status: string(cart.Status),
		Items:  cartItemResponse,
	}
}
