package category

import (
	"bytes"
	"errors"
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
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{},
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
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{
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
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{
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
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{
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
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{
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
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{
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
			name: "categoryService_DeleteCategory_Failed",
			fields: fields{
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{
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
				categoryRepo: &categoryMockRepository{
					Items: []model.Category{
						{
							Base: model.Base{ID: uuid.New()},
						},
					},
				},
			},
			args: args{
				filename: createFile(),
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

func createFile() *bytes.Buffer {
	d1 := []byte("category_name,category_description\ncategory name 1,category description 1\ncategory name 2,category description 2")
	body := bytes.NewBuffer(d1)
	return body
}

type categoryMockRepository struct {
	Items []model.Category
}

func (r *categoryMockRepository) InsertCategory(category *model.Category) error {
	for _, item := range r.Items {
		item.ID = category.ID
		return errors.New("category already exists")
	}
	r.Items = append(r.Items, *category)

	return nil
}

// GetCategories returns all categories
func (r *categoryMockRepository) GetCategories() (*[]model.Category, error) {
	return &r.Items, nil
}

// GetCategoryByID returns a category by id
func (r *categoryMockRepository) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	category := &model.Category{}
	for _, item := range r.Items {
		if item.ID == id {
			category = &item
			return category, nil
		}
	}
	return nil, errors.New("category not found")
}

// UpdateCategory updates a category with the given id
func (r *categoryMockRepository) UpdateCategory(category *model.Category) error {
	for i, item := range r.Items {
		if item.ID == category.ID {
			r.Items[i] = *category
		}
	}
	return nil
}

// InsertBulkCategory inserts bulk categories //TODO
func (r *categoryMockRepository) InsertBulkCategory(categories *[]model.Category) error {
	for _, category := range *categories {
		for _, item := range r.Items {
			if item.Name == category.Name {
				return errors.New("category already exists")
			}
		}
		r.Items = append(r.Items, category)
	}
	return nil
}

// DeleteCategory deletes a category by id
func (r *categoryMockRepository) Delete(category *model.Category) error {
	for index, item := range r.Items {
		if item.ID == category.ID {
			r.Items = append(r.Items[:index], r.Items[index+1:]...)
			return nil
		}
	}
	return errors.New("category not found")
}
