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

func (r *RoleRepository) InsertRole(role *model.Role) error {
	result := r.db.Where("name = ?", role.Name).FirstOrCreate(role)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetRoles returns all roles
func (r *RoleRepository) GetRoles() ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// DeleteRole deletes a role
func (r *RoleRepository) DeleteRole(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&model.Role{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateRole updates a role
func (r *RoleRepository) UpdateRole(role *model.Role) error {
	if err := r.db.Save(role).Error; err != nil {
		return err
	}
	return nil
}

// getRole returns a role
func (r *RoleRepository) GetRole(id string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}
