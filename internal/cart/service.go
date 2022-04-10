package cart

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	product "patika-ecommerce/internal/product"

	"github.com/go-openapi/strfmt"
)

type CartService struct {
	cartRepo     *CartRepository
	cartItemRepo *CartItemRepository
	productRepo  *product.ProductRepository
}

// NewCartService creates a new CartService
func NewCartService(cartRepo *CartRepository, productRepo *product.ProductRepository, cartItemRepo *CartItemRepository) *CartService {
	return &CartService{
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
		productRepo:  productRepo,
	}
}

// GetOrCreateCart returns a cart by user id
func (r *CartService) GetOrCreateCart(user model.User) (*model.Cart, error) {
	return r.cartRepo.GetOrCreateCart(user)
}

// AddToCart adds a product to cart
func (r *CartService) AddToCart(user model.User, req *api.CartAddRequest) (*model.Cart, error) {
	cart, err := r.cartRepo.GetCreatedCart(user)

	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, fmt.Errorf("cart not found")
	}

	product, err := r.productRepo.GetProductWithoutCategories(req.ProductID)
	if err != nil {
		return nil, err
	}

	if *product.Stock < req.Quantity {
		return nil, fmt.Errorf("product stock is not enough")
	}

	is_exists := false
	for _, item := range cart.Items {
		if item.ProductID == req.ProductID {
			item.Quantity += int(req.Quantity)
			r.cartItemRepo.UpdateCartItem(&item)
			is_exists = true
			break
		}
	}
	if !is_exists {
		product, err := r.productRepo.GetProductWithoutCategories(req.ProductID)
		if err != nil {
			return nil, err
		}

		if err := r.cartItemRepo.Create(cart, product); err != nil {
			return nil, err
		}
	}
	cart, err = r.cartRepo.GetCartByID(cart.ID)

	return cart, nil
}

// UpdateCartItem updates a cart item
func (r *CartService) UpdateCartItem(user model.User, id strfmt.UUID, req *api.CartItemUpdateRequest) (*model.CartItem, error) {
	cart, err := r.cartRepo.GetCreatedCartWithItemsAndProducts(user)
	if err != nil {
		return nil, err
	}

	cartItem, err := r.cartItemRepo.GetCartItemByCartAndIDWithProduct(cart, id)
	if err != nil {
		return nil, err
	}

	if req.Quantity > *cartItem.Product.Stock {
		return nil, fmt.Errorf("product stock is not enough")
	}

	cartItem.Quantity = int(req.Quantity)
	if err := r.cartItemRepo.UpdateCartItem(cartItem); err != nil {
		return nil, err
	}

	return cartItem, nil
}

// DeleteCartItem deletes a cart item
func (r *CartService) DeleteCartItem(user model.User, id strfmt.UUID) error {
	cart, err := r.cartRepo.GetCreatedCartWithItemsAndProducts(user) // TODO
	if err != nil {
		return err
	}

	cartItem, err := r.cartItemRepo.GetCartItemByCartAndIDWithProduct(cart, id)
	if err != nil {
		return err
	}

	if err := r.cartItemRepo.DeleteCartItem(cartItem); err != nil {
		return err
	}

	return nil
}
