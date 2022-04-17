package auth

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"patika-ecommerce/internal/model"
)

type AuthRepository interface {
	InsertUser(user *model.User) (*model.User, error)
	GetUser(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
}

type UserRepositoryForMock interface {
	InsertUser(user *model.User) (*model.User, error)
	GetUser(id string) (*model.User, error)
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
	zap.L().Debug("user.repo.insertUser", zap.Reflect("user", user))

	result := u.db.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil

}

// GetUser
func (u *UserRepository) GetUser(id string) (*model.User, error) {
	zap.L().Debug("user.repo.GetUser", zap.Reflect("id", id))

	var user model.User

	result := u.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail
func (u *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	zap.L().Debug("user.repo.GetUserByEmail", zap.Reflect("email", email))

	var user model.User

	result := u.db.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
