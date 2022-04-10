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
