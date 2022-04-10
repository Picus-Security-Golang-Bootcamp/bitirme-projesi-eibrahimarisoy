package cart

import (
	"fmt"
	"patika-ecommerce/internal/model"

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

	// result := r.db.Where("name = ?", category.Name).FirstOrCreate(category)

	return r.db.Create(cartItem).Error
}
