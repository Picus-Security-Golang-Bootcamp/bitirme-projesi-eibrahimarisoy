package category

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
)

// CategoryRequestToCategory converts a CategoryRequest to a Category
func CategoryRequestToCategory(categoryRequest *api.CategoryRequest) *model.Category {
	return &model.Category{
		Name:        *categoryRequest.Name,
		Description: categoryRequest.Description,
	}
}

// CategoryToCategoryResponse converts a Category to a CategoryResponse
func CategoryToCategoryResponse(category *model.Category) *api.CategoryResponse {
	return &api.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
	}
}
