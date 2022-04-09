package product

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"

	"github.com/go-openapi/strfmt"
)

func ProductRequestToProduct(productRequest *api.ProductRequest) *model.Product {
	stockAddr := productRequest.Stock
	stock := int(*stockAddr)

	categories := []strfmt.UUID{}
	for _, category := range productRequest.Categories {
		categories = append(categories, category)
	}

	return &model.Product{
		Name:         productRequest.Name,
		Description:  *productRequest.Description,
		Price:        productRequest.Price,
		Stock:        &stock,
		SKU:          productRequest.Sku,
		CategoriesID: categories,
	}
}
