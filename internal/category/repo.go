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

// InsertCategory inserts a new category
func (r *CategoryRepository) InsertCategory(category *model.Category) error {
	return r.db.Save(category).Error
}

// GetCategories returns all categories
func (r *CategoryRepository) GetCategories() ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
