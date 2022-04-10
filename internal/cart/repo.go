package cart

import (
	"fmt"
	"patika-ecommerce/internal/model"

	"github.com/go-openapi/strfmt"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

type CartItemRepository struct {
	db *gorm.DB
}

func (r *CartRepository) Migration() {
	r.db.AutoMigrate(&model.Cart{})
}

func (r *CartItemRepository) Migration() {
	r.db.AutoMigrate(&model.CartItem{})
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func NewCartItemRepository(db *gorm.DB) *CartItemRepository {
	return &CartItemRepository{db: db}
}

// GetOrCreateCart returns a cart by user id
func (r *CartRepository) GetOrCreateCart(user model.User) (*model.Cart, error) {
	cart := &model.Cart{}
	if err := r.db.Preload("Items.Product").Where("user_id = ?", user.ID).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cart.UserID = user.ID
			if err := r.db.Create(cart).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	fmt.Println("Created cart", cart)

	return cart, nil
}

// GetCreatedCart returns a cart by user id
func (r *CartRepository) GetCreatedCart(user model.User) (*model.Cart, error) {
	cart := &model.Cart{}
	if err := r.db.Preload("Items").Where("user_id = ? AND status = ?", user.ID, model.CartStatusCreated).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	fmt.Println("Created cart", cart)
	return cart, nil
}

// GetCreatedCartWithItemsAndProducts returns a cart by user id
func (r *CartRepository) GetCreatedCartWithItemsAndProducts(user model.User) (*model.Cart, error) {
	cart := &model.Cart{}
	if err := r.db.Preload("Items.Product").Where("user_id = ? AND status = ?", user.ID, model.CartStatusCreated).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return cart, nil
}

// GetCartByID returns a cart by id
func (r *CartRepository) GetCartByID(id strfmt.UUID) (*model.Cart, error) {
	cart := &model.Cart{}
	if err := r.db.Preload("Items.Product").Where("id = ?", id).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	fmt.Println("Created cart", cart)
	return cart, nil
}

// UpdateCart updates a cart
func (r *CartRepository) UpdateCart(cart *model.Cart) error {

	return r.db.Model(&cart).Updates(cart).Error
}

// ###### CART ITEM ######

func (r *CartItemRepository) Create(cart *model.Cart, product *model.Product) error {
	cartItem := &model.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  1,
		Price:     *product.Price,
	}
	// cart.Items = append(cart.Items, *cartItem)

	return r.db.Create(cartItem).Error
}

// UpdateCartItem updates a cart item
func (r *CartItemRepository) UpdateCartItem(cartItem *model.CartItem) error {
	return r.db.Model(&cartItem).Updates(cartItem).Error
}

// GetCartItemByID returns a cart item by id
func (r *CartItemRepository) GetCartItemByCartAndIDWithProduct(cart *model.Cart, id strfmt.UUID) (*model.CartItem, error) {
	cartItem := &model.CartItem{}
	if err := r.db.Model(&cartItem).Preload("Product").Where("cart_id = ? AND id = ?", cart.ID, id).First(cartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return cartItem, nil
}

// DeleteCartItem deletes a cart item
func (r *CartItemRepository) DeleteCartItem(cartItem *model.CartItem) error {
	return r.db.Delete(cartItem).Error
}
