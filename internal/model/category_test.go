package model

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func TestCategory_BeforeCreate(t *testing.T) {
	name := "test"
	description := "test"
	type fields struct {
		Base        Base
		Name        *string
		Slug        string
		Description string
		Products    []Product
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
					ID: uuid.New(),
				},
				Name:        &name,
				Description: description,
				Products:    []Product{},
			},
			args: args{
				tx: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Category{
				Base:        tt.fields.Base,
				Name:        tt.fields.Name,
				Slug:        tt.fields.Slug,
				Description: tt.fields.Description,
				Products:    tt.fields.Products,
			}
			if err := c.BeforeCreate(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("Category.BeforeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, c.Slug, slug.Make(name))
		})
	}
}
