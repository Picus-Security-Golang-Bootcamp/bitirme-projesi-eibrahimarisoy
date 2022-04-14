package product

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"

	"github.com/go-openapi/strfmt"
)

func ProductRequestToProduct(productRequest *api.ProductRequest) *model.Product {
	stockAddr := productRequest.Stock
	stock := int64(*stockAddr)

	categories := []model.Category{}
	// category_ids := []strfmt.UUID{}
	for _, c := range productRequest.Categories {
		categories = append(categories, model.Category{Base: model.Base{ID: c.ID}})
		// category_ids = append(category_ids, c.ID)
	}

	return &model.Product{
		Name:        productRequest.Name,
		Description: *productRequest.Description,
		Price:       productRequest.Price,
		Stock:       &stock,
		SKU:         productRequest.Sku,
		Categories:  categories,
		// CategoriesID: category_ids,
	}
}

func ProductToResponse(product *model.Product) *api.ProductResponse {
	stock := int64(*product.Stock)
	categories := []strfmt.UUID{}
	for _, c := range product.Categories {
		categories = append(categories, c.ID)
	}

	return &api.ProductResponse{
		ID:          product.ID,
		Slug:        product.Slug,
		Name:        *product.Name,
		Description: product.Description,
		Price:       *product.Price,
		Stock:       stock,
		Sku:         *product.SKU,
		Categories:  categories,
	}
}

func ProductsToResponse(products *[]model.Product) []*api.ProductResponse {
	response := []*api.ProductResponse{}
	for _, product := range *products {
		response = append(response, ProductToResponse(&product))
	}
	return response
}

func ProductToProductBasicResponse(product *model.Product) *api.ProductBasicResponse {
	stock := int64(*product.Stock)

	return &api.ProductBasicResponse{
		ID:          product.ID,
		Slug:        product.Slug,
		Name:        *product.Name,
		Description: product.Description,
		Price:       *product.Price,
		Stock:       stock,
	}
}

// func ProductUpdateRequestToProduct(productUpdateRequest *api.ProductUpdateRequest) *model.Product {
// 	stockAddr := productUpdateRequest.Stock
// 	stock := int(stockAddr)

// 	categories := []model.Category{}
// 	for _, c := range productUpdateRequest.Categories {
// 		categories = append(categories, model.Category{Base: model.Base{ID: c.ID}})
// 	}

// 	fmt.Printf("%+v", productUpdateRequest)

// 	a := &model.Product{
// 		Name:        &productUpdateRequest.Name,
// 		Description: productUpdateRequest.Description,
// 		Price:       &productUpdateRequest.Price,
// 		Stock:       &stock,
// 		Categories:  &categories,
// 	}
// 	fmt.Println("ddddddddddd", a.ToString())

// 	return a
// }
