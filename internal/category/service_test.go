package category

import (
	"bytes"
	"errors"
	"fmt"
	"patika-ecommerce/internal/model"
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func TestCategoryService_CreateCategory(t *testing.T) {
	categoryName := "category name 1"
	categoryDescription := "category description 1"

	type fields struct {
		categoryRepo MockCategoryRepository
	}
	type args struct {
		category *model.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "categoryService_CreateCategory_Success",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{},
				},
			},
			args: args{
				category: &model.Category{
					Base:        model.Base{ID: uuid.New()},
					Name:        &categoryName,
					Description: categoryDescription,
				},
			},
			wantErr: false,
		},
		{
			name: "categoryService_CreateCategory_Failed_AlreadyExists",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{
						{
							Base:        model.Base{ID: uuid.New()},
							Name:        &categoryName,
							Description: categoryDescription,
						},
					},
				},
			},
			args: args{
				category: &model.Category{
					Base:        model.Base{ID: uuid.New()},
					Name:        &categoryName,
					Description: categoryDescription,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryService{
				categoryRepo: tt.fields.categoryRepo,
			}
			if err := c.CreateCategory(tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("CategoryService.CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryService_GetCategories(t *testing.T) {
	categoryName := "category name 1"
	categoryDescription := "category description 1"

	type fields struct {
		categoryRepo MockCategoryRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    *[]model.Category
		wantErr bool
	}{
		{
			name: "categoryService_GetCategories_Success",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{
						{
							Name:        &categoryName,
							Description: categoryDescription,
						},
					},
				},
			},
			want: &[]model.Category{
				{
					Name:        &categoryName,
					Description: categoryDescription,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryService{
				categoryRepo: tt.fields.categoryRepo,
			}
			got, err := c.GetCategories()
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryService.GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryService.GetCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryService_GetCategoryByID(t *testing.T) {
	categoryName := "category name 1"
	categoryDescription := "category description 1"
	id := uuid.New()

	type fields struct {
		categoryRepo MockCategoryRepository
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Category
		wantErr bool
	}{
		{
			name: "categoryService_GetCategoryByID_Success",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{
						{
							Base:        model.Base{ID: id},
							Name:        &categoryName,
							Description: categoryDescription,
						},
					},
				},
			},
			args: args{
				id: id,
			},
			want: &model.Category{
				Base:        model.Base{ID: id},
				Name:        &categoryName,
				Description: categoryDescription,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryService{
				categoryRepo: tt.fields.categoryRepo,
			}
			got, err := c.GetCategoryByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryService.GetCategoryByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryService.GetCategoryByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	categoryName := "category name 1"
	categoryDescription := "category description 1"
	id := uuid.New()

	newCategoryName := "new category name"
	newCategoryDescription := "new category description"

	type fields struct {
		categoryRepo MockCategoryRepository
	}
	type args struct {
		category *model.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "categoryService_GetCategoryByID_Success",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{
						{
							Base:        model.Base{ID: id},
							Name:        &categoryName,
							Description: categoryDescription,
						},
					},
				},
			},
			args: args{
				category: &model.Category{
					Base:        model.Base{ID: id},
					Name:        &newCategoryName,
					Description: newCategoryDescription,
				},
			},
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryService{
				categoryRepo: tt.fields.categoryRepo,
			}
			if err := c.UpdateCategory(tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("CategoryService.UpdateCategory() error = %v, wantErr %v", err, tt.wantErr)
				assert.Equal(t, tt.args.category.Name, newCategoryName)
			}
		})
	}
}

func TestCategoryService_DeleteCategory(t *testing.T) {
	id := uuid.New()
	idTwo := uuid.New()
	type fields struct {
		categoryRepo MockCategoryRepository
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "categoryService_DeleteCategory_Success",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{
						{
							Base: model.Base{ID: id},
						},
					},
				},
			},
			args: args{
				id: id,
			},
			wantErr: false,
		},
		{
			name: "categoryService_DeleteCategory_Success",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{
						{
							Base: model.Base{ID: id},
						},
					},
				},
			},
			args: args{
				id: idTwo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryService{
				categoryRepo: tt.fields.categoryRepo,
			}
			if err := c.DeleteCategoryService(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("CategoryService.DeleteCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCategoryService_CreateBulkCategories(t *testing.T) {
	fmt.Println("TestCategoryService_CreateBulkCategories")

	type fields struct {
		categoryRepo MockCategoryRepository
	}
	type args struct {
		filename *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Category
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "categoryService_CreateBulkCategories_Success",
			fields: fields{
				categoryRepo: &mockCategoryRepository{
					items: []model.Category{
						{
							Base: model.Base{ID: uuid.New()},
						},
					},
				},
			},
			args: args{
				filename: createFilename(),
			},
			want: []model.Category{
				{
					Base: model.Base{ID: uuid.New()},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryService{
				categoryRepo: tt.fields.categoryRepo,
			}
			_, err := c.CreateBulkCategories(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryService.CreateBulkCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func createFilename() *bytes.Buffer {
	d1 := []byte("category_name,category_description\ncategory name 222,category description 1\ncategory name 3332,category description     2")
	body := bytes.NewBuffer(d1)
	return body
}

type mockCategoryRepository struct {
	items []model.Category
}

func (r *mockCategoryRepository) InsertCategory(category *model.Category) error {
	for _, item := range r.items {
		item.ID = category.ID
		return errors.New("category already exists")
	}
	r.items = append(r.items, *category)

	return nil
}

// GetCategories returns all categories
func (r *mockCategoryRepository) GetCategories() (*[]model.Category, error) {
	return &r.items, nil
}

// GetCategoryByID returns a category by id
func (r *mockCategoryRepository) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	category := &model.Category{}
	for _, item := range r.items {
		if item.ID == id {
			category = &item
		}
	}
	return category, nil
}

// UpdateCategory updates a category with the given id
func (r *mockCategoryRepository) UpdateCategory(category *model.Category) error {
	for i, item := range r.items {
		if item.ID == category.ID {
			r.items[i] = *category
		}
	}

	return nil
}

// InsertBulkCategory inserts bulk categories //TODO
func (r *mockCategoryRepository) InsertBulkCategory(categories *[]model.Category) error {
	for _, category := range *categories {
		for _, item := range r.items {
			if item.Name == category.Name {
				return errors.New("category already exists")
			}
		}
		r.items = append(r.items, category)
	}

	return nil
}

// DeleteCategory deletes a category by id
func (r *mockCategoryRepository) Delete(category *model.Category) error {
	isExist := false
	for index, item := range r.items {
		if item.ID == category.ID {
			category = &item
			isExist = true
			r.items = append(r.items[:index], r.items[index+1:]...)
		}
	}
	if !isExist {
		return errors.New("category not found")
	}
	return nil
}
