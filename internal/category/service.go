package category

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"patika-ecommerce/internal/model"
)

type CategoryService struct {
	categoryRepo *CategoryRepository
}

func NewCategoryService(categoryRepo *CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
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
			Name:        record[0],
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
