package category

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"patika-ecommerce/internal/model"

	"github.com/go-openapi/strfmt"
)

type CategoryService struct {
	categoryRepo *CategoryRepository
}

func NewCategoryService(categoryRepo *CategoryRepository) *CategoryService {
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
func (c *CategoryService) GetCategoryByID(id strfmt.UUID4) (*model.Category, error) {
	return c.categoryRepo.GetCategoryByID(id)
}

// UpdateCategory updates a category
func (c *CategoryService) UpdateCategory(category *model.Category) error {
	return c.categoryRepo.UpdateCategory(category)
}

func (c *CategoryService) CreateBulkCategories(filename *multipart.FileHeader) ([]model.Category, error) {
	file, err := filename.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var categories []model.Category
	for _, record := range records[1:] {
		category := model.Category{
			Name:        &record[0],
			Description: record[1],
			// ParentId:    record[2],
		}
		categories = append(categories, category)
	}
	fmt.Println(c.categoryRepo)
	for _, category := range categories {

		if err := c.categoryRepo.InsertCategory(&category); err != nil {
			return nil, err
		}
	}
	return categories, nil

}
