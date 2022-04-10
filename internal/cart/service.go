package cart

import "patika-ecommerce/internal/model"

type CartService struct {
	cartRepo *CartRepository
}

// NewCartService creates a new CartService
func NewCartService(cartRepo *CartRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}

// GetOrCreateCart returns a cart by user id
func (r *CartService) GetOrCreateCart(user model.User) (*model.Cart, error) {
	return r.cartRepo.GetOrCreateCart(user)
}
