package category

import (
	"patika-ecommerce/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func (r *CategoryRepository) Migration() {
	r.db.AutoMigrate(&model.Category{})
}
func NewCategoryrRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}
