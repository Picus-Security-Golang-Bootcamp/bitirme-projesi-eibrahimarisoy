package cart

import (
	"patika-ecommerce/internal/model"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func (r *CartRepository) Migration() {
	r.db.AutoMigrate(&model.Cart{})
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// GetOrCreateCart returns a cart by user id
func (r *CartRepository) GetOrCreateCart(user model.User) (*model.Cart, error) {
	cart := &model.Cart{}
	if err := r.db.Where("user_id = ?", user.ID).First(cart).Error; err != nil {
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
