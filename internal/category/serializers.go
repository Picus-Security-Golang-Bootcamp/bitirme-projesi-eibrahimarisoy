package category

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	common "patika-ecommerce/pkg/utils"
)

// CategoryRequestToCategory converts a CategoryRequest to a Category
func CategoryRequestToCategory(categoryRequest *api.CategoryRequest) *model.Category {
	return &model.Category{
		Name:        categoryRequest.Name,
		Description: categoryRequest.Description,
	}
}

// CategoryToCategoryResponse converts a Category to a CategoryResponse
func CategoryToCategoryResponse(category *model.Category) *api.CategoryResponse {
	fmt.Println(category)
	return &api.CategoryResponse{
		ID:          common.UUIDToStrfmt(category.ID),
		Name:        *category.Name,
		Slug:        category.Slug,
		Description: category.Description,
	}
}

// CategoriesToCategoryResponse converts a list of Categories to a list of CategoryResponse
func CategoriesToCategoryResponse(categories *[]model.Category) []*api.CategoryResponse {
	var categoryResponses []*api.CategoryResponse
	for _, category := range *categories {
		categoryResponses = append(categoryResponses, CategoryToCategoryResponse(&category))
	}
	return categoryResponses
}
