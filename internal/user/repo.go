package auth

import (
	"gorm.io/gorm"

	"patika-ecommerce/internal/model"
)

type AuthRepository interface {
	InsertUser(user *model.User) (*model.User, error)
	GetAll() (*[]model.User, error)
	GetUser(id string) (*model.User, error)
	DeleteUser(id string) error
	UpdateUser(user *model.User) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}

type UserRepositoryForMock interface {
	InsertUser(user *model.User) (*model.User, error)
	GetAll() (*[]model.User, error)
	GetUser(id string) (*model.User, error)
	DeleteUser(id string) error
	UpdateUser(user *model.User) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}
type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Migration() {
	r.db.AutoMigrate(&model.User{})
}
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// InsertUser insert user to database
func (u *UserRepository) InsertUser(user *model.User) (*model.User, error) {
	result := u.db.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil

}

func (u *UserRepository) GetAll() (*[]model.User, error) {
	var users []model.User

	result := u.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

// GetUser
func (u *UserRepository) GetUser(id string) (*model.User, error) {
	var user model.User

	result := u.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// DeleteUser
func (u *UserRepository) DeleteUser(id string) error {
	result := u.db.Debug().Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateUser
func (u *UserRepository) UpdateUser(user *model.User) (*model.User, error) {
	result := u.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// GetUserByEmail
func (u *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	result := u.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
