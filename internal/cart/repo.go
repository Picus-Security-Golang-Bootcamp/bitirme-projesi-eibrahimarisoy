package cart

import (
	"errors"
	"patika-ecommerce/internal/model"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MockCartRepository interface {
	GetOrCreateCart(user *model.User) (*model.Cart, error)
	GetCreatedCart(user *model.User) (*model.Cart, error)
	GetCreatedCartWithItemsAndProducts(user *model.User) (*model.Cart, error)
	GetCreatedCartWithItems(user *model.User) (*model.Cart, error)
	GetCreatedCartByUserAndCart(user *model.User, cartId strfmt.UUID) (*model.Cart, error)
	GetCartByID(id uuid.UUID) (*model.Cart, error)
	UpdateCart(cart *model.Cart) error
}

type MockCartItemRepository interface {
	Create(cart *model.Cart, product *model.Product) error
	UpdateCartItem(cartItem *model.CartItem) error
	GetCartItemByCartAndIDWithProduct(cart *model.Cart, id uuid.UUID) (*model.CartItem, error)
	DeleteCartItem(cartItem *model.CartItem) error
}

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

// GetOrCreateCart if cart is exists returns it otherwise create cart and return it
func (r *CartRepository) GetOrCreateCart(user *model.User) (*model.Cart, error) {
	zap.L().Debug("cart.repo.GetOrCreateCart", zap.Reflect("user", user))

	cart := &model.Cart{}
	result := r.db.Model(&user).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if err := r.db.Preload("Items.Product").Where("user_id = ? AND status = ?", user.ID, model.CartStatusCreated).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cart.UserID = user.ID
			if err := r.db.Create(cart).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return cart, nil
}

// GetCreatedCart returns a created cart by user id
func (r *CartRepository) GetCreatedCart(user *model.User) (*model.Cart, error) {
	zap.L().Debug("cart.repo.GetCreatedCart", zap.Reflect("user", user))

	cart := &model.Cart{}
	if err := r.db.Preload("Items").Where("user_id = ? AND status = ?", user.ID, model.CartStatusCreated).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Cart not found. Please create a cart")
		}
		return nil, err
	}
	return cart, nil
}

// GetCreatedCartWithItemsAndProducts returns a cart by user id
func (r *CartRepository) GetCreatedCartWithItemsAndProducts(user *model.User) (*model.Cart, error) {
	zap.L().Debug("cart.repo.GetCreatedCartWithItemsAndProducts", zap.Reflect("user", user))

	cart := &model.Cart{}
	if err := r.db.Preload("Items.Product").Where("user_id = ? AND status = ?", user.ID, model.CartStatusCreated).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Cart not found. Please create a cart")
		}
		return nil, err
	}
	return cart, nil
}

// GetCreatedCartWithItems returns a cart by user id
func (r *CartRepository) GetCreatedCartWithItems(user *model.User) (*model.Cart, error) {
	zap.L().Debug("cart.repo.GetCreatedCartWithItems", zap.Reflect("user", user))

	cart := &model.Cart{}
	if err := r.db.Preload("Items").Where("user_id = ? AND status = ?", user.ID, model.CartStatusCreated).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("Cart not found. Please create a cart")
		}
		return nil, err
	}
	return cart, nil
}

// GetCreatedCartByUserAndCart returns a cart by user id
func (r *CartRepository) GetCreatedCartByUserAndCart(user *model.User, cartId strfmt.UUID) (*model.Cart, error) {
	zap.L().Debug("cart.repo.GetCreatedCartByUserAndCart", zap.Reflect("user", user), zap.Reflect("cartId", cartId))

	cart := model.Cart{}
	if err := r.db.Preload("Items.Product").
		Where("user_id = ? AND status = ? AND id = ?", user.ID, model.CartStatusCreated, cartId).
		First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

// GetCartByID returns a cart by id
func (r *CartRepository) GetCartByID(id uuid.UUID) (*model.Cart, error) {
	zap.L().Debug("cart.repo.GetCartByID", zap.Reflect("id", id))

	cart := &model.Cart{}
	if err := r.db.Preload("Items.Product").Where("id = ?", id).First(cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return cart, nil
}

// UpdateCart updates a cart
func (r *CartRepository) UpdateCart(cart *model.Cart) error {
	zap.L().Debug("cart.repo.UpdateCart", zap.Reflect("cart", cart))

	return r.db.Model(&cart).Updates(cart).Error
}

// ###### CART ITEM ######

func (r *CartItemRepository) Create(cart *model.Cart, product *model.Product) error {
	zap.L().Debug("cartItem.repo.Create", zap.Reflect("cart", cart), zap.Reflect("product", product))

	cartItem := &model.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  1,
		Price:     product.Price,
	}
	cart.Items = append(cart.Items, *cartItem)

	return r.db.Create(cartItem).Error
}

// UpdateCartItem updates a cart item
func (r *CartItemRepository) UpdateCartItem(cartItem *model.CartItem) error {
	zap.L().Debug("cartItem.repo.UpdateCartItem", zap.Reflect("cartItem", cartItem))

	return r.db.Model(&cartItem).Updates(cartItem).Error
}

// GetCartItemByID returns a cart item by id
func (r *CartItemRepository) GetCartItemByCartAndIDWithProduct(cart *model.Cart, id uuid.UUID) (*model.CartItem, error) {
	zap.L().Debug("cartItem.repo.GetCartItemByCartAndIDWithProduct", zap.Reflect("cart", cart), zap.Reflect("id", id))

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
	zap.L().Debug("cartItem.repo.DeleteCartItem", zap.Reflect("cartItem", cartItem))

	return r.db.Delete(cartItem).Error
}
