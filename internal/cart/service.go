package cart

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	product "patika-ecommerce/internal/product"
	common "patika-ecommerce/pkg/utils"

	"github.com/google/uuid"
)

type MockCartService interface {
	GetOrCreateCart(user *model.User) (*model.Cart, error)
	AddToCart(user *model.User, req *api.AddToCartRequest) (*model.Cart, error)
	UpdateCartItem(user *model.User, id uuid.UUID, req *api.CartItemUpdateRequest) (*model.CartItem, error)
	DeleteCartItem(user *model.User, id uuid.UUID) error
}

type CartService struct {
	cartRepo     MockCartRepository
	cartItemRepo MockCartItemRepository
	productRepo  product.MockProductRepository
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

	pId, err := common.StrfmtToUUID(req.ProductID)
	if err != nil {
		return nil, err
	}

	// find product by given id
	product, err := r.productRepo.GetProductWithoutCategories(pId)
	if err != nil {
		return nil, err
	}

	// check if product is available in cart
	is_exists := false
	for _, item := range cart.Items {
		// if product is already in cart then update quantity
		if item.ProductID == pId {
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
	return cart, nil
}

// UpdateCartItem updates a cart item
func (r *CartService) UpdateCartItem(user *model.User, id uuid.UUID, req *api.CartItemUpdateRequest) (*model.CartItem, error) {
	cart, err := r.cartRepo.GetCreatedCartWithItemsAndProducts(user)
	if err != nil {
		return nil, err
	}

	cartItem, err := cart.GetCartItemByID(id)
	if err != nil {
		return nil, err
	}

	if req.Quantity == 0 {
		r.cartItemRepo.DeleteCartItem(cartItem)
		cartItem.Quantity = 0
		return cartItem, nil
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
func (r *CartService) DeleteCartItem(user *model.User, id uuid.UUID) error {
	cart, err := r.cartRepo.GetCreatedCartWithItems(user)
	if err != nil {
		return err
	}

	cartItem, err := cart.GetCartItemByID(id)
	if err != nil {
		return err
	}

	if err := r.cartItemRepo.DeleteCartItem(cartItem); err != nil {
		return err
	}

	return nil
}
