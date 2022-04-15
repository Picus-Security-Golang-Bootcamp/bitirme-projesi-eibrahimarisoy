package category

import (
	"errors"
	"patika-ecommerce/internal/model"
	"testing"

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
	categories := &[]model.Category{}
	// //  &[]models.Book{}
	// if err := r.db.Find(categories).Error; err != nil {
	// 	return nil, err
	// }

	return categories, nil
}

// GetCategoryByID returns a category by id
func (r *mockCategoryRepository) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	category := &model.Category{}
	// if err := r.db.Where("id = ?", id).First(category).Error; err != nil {
	// 	return nil, err
	// }
	return category, nil
}

// UpdateCategory updates a category with the given id
func (r *mockCategoryRepository) UpdateCategory(category *model.Category) error {
	// result := r.db.Updates(category)
	// if result.Error != nil {
	// 	return result.Error
	// }
	return nil
}

// InsertBulkCategory inserts bulk categories
func (r *mockCategoryRepository) InsertBulkCategory(categories *[]model.Category) error {
	// tx := r.db.Begin()

	// for _, category := range *categories {
	// 	if err := tx.Create(&category).Error; err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }

	// tx.Commit()
	return nil
}

// DeleteCategory deletes a category by id
func (r *mockCategoryRepository) DeleteCategory(category *model.Category) error {
	// result := r.db.Select(clause.Associations).Delete(category)
	// if result.Error != nil {
	// 	return result.Error
	// }
	return nil
}
