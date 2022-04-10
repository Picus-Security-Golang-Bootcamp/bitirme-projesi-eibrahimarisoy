package cart

type CartService struct {
	cartRepo *CartRepository
}

// NewCartService creates a new CartService
func NewCartService(cartRepo *CartRepository) *CartService {
	return &CartService{
		cartRepo: cartRepo,
	}
}
