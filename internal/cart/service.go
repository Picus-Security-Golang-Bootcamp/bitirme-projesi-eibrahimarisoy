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
func (r *CartService) GetOrCreateCart(user *model.User) (*model.Cart, error) {
	return r.cartRepo.GetOrCreateCart(user)
}

// AddToCart adds a product to cart
func (r *CartService) AddToCart(user *model.User, req *api.AddToCartRequest) (*model.Cart, error) {
	cart, err := r.cartRepo.GetCreatedCart(user)
	if err != nil {
		return nil, err
	}

	// find product by given id
	product, err := r.productRepo.GetProductWithoutCategories(req.ProductID)
	if err != nil {
		return nil, err
	}

	// check if product is available in cart
	is_exists := false
	for _, item := range cart.Items {
		// if product is already in cart then update quantity
		if item.ProductID == req.ProductID {
			if int64(item.Quantity)+req.Quantity > *product.Stock {
				return nil, fmt.Errorf("Product stock is not enough")
			}
			item.Quantity += req.Quantity
			r.cartItemRepo.UpdateCartItem(&item)
			is_exists = true
			break
		}
	}

	// if product not exists in cart, create new cart item
	if !is_exists {
		if *product.Stock < req.Quantity {
			return nil, fmt.Errorf("Product stock is not enough")
		}

		if err := r.cartItemRepo.Create(cart, product); err != nil {
			return nil, err
		}
	}

	return r.cartRepo.GetCartByID(cart.ID)
}

// UpdateCartItem updates a cart item
func (r *CartService) UpdateCartItem(user *model.User, id strfmt.UUID, req *api.CartItemUpdateRequest) (*model.CartItem, error) {
	cart, err := r.cartRepo.GetCreatedCartWithItemsAndProducts(user)
	if err != nil {
		return nil, err
	}

	cartItem, err := cart.GetCartItemByID(id)
	if err != nil {
		return nil, err
	}

	if req.Quantity > *cartItem.Product.Stock {
		return nil, fmt.Errorf("product stock is not enough")
	}

	cartItem.Quantity = req.Quantity
	if err := r.cartItemRepo.UpdateCartItem(cartItem); err != nil {
		return nil, err
	}

	return cartItem, nil
}

// DeleteCartItem deletes a cart item
func (r *CartService) DeleteCartItem(user *model.User, id strfmt.UUID) error {
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
