package category

import (
	"bytes"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/utils"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepo CategoryRepositoryInterface
}

type MockCategoryService interface {
	CreateCategory(category *model.Category) error
	GetCategories() (*[]model.Category, error)
	GetCategoryByID(id uuid.UUID) (*model.Category, error)
	UpdateCategory(category *model.Category) error
	DeleteCategoryService(id uuid.UUID) error
	CreateBulkCategories(filename *bytes.Buffer) ([]model.Category, error)
}

func NewCategoryService(categoryRepo CategoryRepositoryInterface) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategory creates a new category
func (c *CategoryService) CreateCategory(category *model.Category) error {
	return c.categoryRepo.InsertCategory(category)
}

// GetCategories returns all categories
func (c *CategoryService) GetCategories() (*[]model.Category, error) {
	return c.categoryRepo.GetCategories()
}

// GetCategoryByID returns a category by id
func (c *CategoryService) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	return c.categoryRepo.GetCategoryByID(id)
}

// UpdateCategory updates a category
func (c *CategoryService) UpdateCategory(category *model.Category) error {
	return c.categoryRepo.UpdateCategory(category)
}

// CreateBulkCategories creates multiple categories in bulk operation with the specified file
func (c *CategoryService) CreateBulkCategories(buf *bytes.Buffer) ([]model.Category, error) {
	records, err := utils.ReadFile(buf)
	if err != nil {
		return nil, err
	}

	var categories []model.Category
	for _, record := range records[1:] {
		category := model.Category{
			Name:        &record[0],
			Description: record[1],
		}
		categories = append(categories, category)
	}

	if err := c.categoryRepo.InsertBulkCategory(&categories); err != nil {
		return nil, err
	}
	return categories, nil

}

// DeleteCategory deletes a category by id
func (c *CategoryService) DeleteCategoryService(id uuid.UUID) error {
	category, err := c.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return err
	}
	return c.categoryRepo.Delete(category)
}
