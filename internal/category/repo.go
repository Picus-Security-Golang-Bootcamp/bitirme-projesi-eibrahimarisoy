package category

import (
	"fmt"
	"patika-ecommerce/internal/model"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryRepository struct {
	db *gorm.DB
}

type CategoryRepositoryInterface interface {
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
	zap.L().Debug("category.repo.InsertCategory", zap.Reflect("category", category))

	result := r.db.Create(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetCategories returns all categories
func (r *CategoryRepository) GetCategories() (*[]model.Category, error) {
	zap.L().Debug("category.repo.GetCategories")

	categories := &[]model.Category{}

	if err := r.db.Find(categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryByID returns a category by id
func (r *CategoryRepository) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	zap.L().Debug("category.repo.GetCategoryByID", zap.Reflect("id", id))

	category := &model.Category{}
	if err := r.db.Where("id = ?", id).First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateCategory updates a category with the given id
func (r *CategoryRepository) UpdateCategory(category *model.Category) error {
	zap.L().Debug("category.repo.UpdateCategory", zap.Reflect("category", category))

	result := r.db.Updates(category)
	if result.Error != nil {
		return result.Error
	}
	fmt.Println(result)
	return nil
}

// InsertBulkCategory inserts bulk categories
func (r *CategoryRepository) InsertBulkCategory(categories *[]model.Category) error {
	zap.L().Debug("category.repo.InsertBulkCategory", zap.Reflect("categories", categories))

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
	zap.L().Debug("category.repo.Delete", zap.Reflect("category", category))

	result := r.db.Debug().Select(clause.Associations).Delete(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
