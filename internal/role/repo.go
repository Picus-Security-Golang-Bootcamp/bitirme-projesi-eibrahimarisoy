package role

import (
	"patika-ecommerce/internal/model"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func (r *RoleRepository) Migration() {
	r.db.AutoMigrate(&model.Role{})
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func InsertRole(db *gorm.DB, role *model.Role) error {
	result := db.Where("name = ?", role.Name).FirstOrCreate(role)
	if result.Error != nil {
		return result.Error
	}
	return nil

}
