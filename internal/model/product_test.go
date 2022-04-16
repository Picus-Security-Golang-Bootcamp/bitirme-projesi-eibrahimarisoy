package model

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestProduct_BeforeCreate(t *testing.T) {
	name := "test"
	stock := int64(1)
	sku := "test"

	type fields struct {
		Base         Base
		Name         *string
		Slug         string
		Description  string
		Price        float64
		Stock        *int64
		SKU          *string
		Categories   []Category
		CategoriesID []strfmt.UUID
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "BeforeCreate_Succeed",
			fields: fields{
				Base: Base{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
				},
				Name:        &name,
				Description: "description",
				Price:       100,
				Stock:       &stock,
				SKU:         &sku,
				Categories: []Category{
					{
						Base: Base{ID: uuid.New()},
					},
				},
			},
			args: args{
				tx: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{
				Base:         tt.fields.Base,
				Name:         tt.fields.Name,
				Slug:         tt.fields.Slug,
				Description:  tt.fields.Description,
				Price:        tt.fields.Price,
				Stock:        tt.fields.Stock,
				SKU:          tt.fields.SKU,
				Categories:   tt.fields.Categories,
				CategoriesID: tt.fields.CategoriesID,
			}
			if err := p.BeforeCreate(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("Product.BeforeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
