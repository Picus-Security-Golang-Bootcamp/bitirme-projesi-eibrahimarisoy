package product

import (
	"patika-ecommerce/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func (r *ProductRepository) Migration() {
	r.db.AutoMigrate(&model.Role{})
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) InsertRole(role *model.Role) error {
	result := r.db.Where("name = ?", role.Name).FirstOrCreate(role)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
