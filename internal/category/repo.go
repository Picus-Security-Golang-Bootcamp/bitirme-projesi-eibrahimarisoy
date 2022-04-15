package category

import (
	"patika-ecommerce/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryRepository struct {
	db *gorm.DB
}

type MockCategoryRepository interface {
	InsertCategory(category *model.Category) error
	GetCategories() (*[]model.Category, error)
	GetCategoryByID(id uuid.UUID) (*model.Category, error)
	UpdateCategory(category *model.Category) error
	InsertBulkCategory(categories *[]model.Category) error
	Delete(category *model.Category) error
}

func (r *CategoryRepository) Migration() {
	r.db.AutoMigrate(&model.Category{})
}
func NewCategoryrRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// InsertCategory inserts a new category
func (r *CategoryRepository) InsertCategory(category *model.Category) error {
	result := r.db.Create(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetCategories returns all categories
func (r *CategoryRepository) GetCategories() (*[]model.Category, error) {
	categories := &[]model.Category{}
	//  &[]models.Book{}
	if err := r.db.Find(categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryByID returns a category by id
func (r *CategoryRepository) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	category := &model.Category{}
	if err := r.db.Where("id = ?", id).First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateCategory updates a category with the given id
func (r *CategoryRepository) UpdateCategory(category *model.Category) error {
	result := r.db.Updates(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// InsertBulkCategory inserts bulk categories
func (r *CategoryRepository) InsertBulkCategory(categories *[]model.Category) error {
	tx := r.db.Begin()

	for _, category := range *categories {
		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// DeleteCategory deletes a category by id
func (r *CategoryRepository) Delete(category *model.Category) error {
	result := r.db.Select(clause.Associations).Delete(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
