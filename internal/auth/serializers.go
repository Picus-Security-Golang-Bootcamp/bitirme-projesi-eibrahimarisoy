package auth

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
)

// RegisterToUser converts a RegisterUser to a User
func RegisterToUser(user *api.RegisterUser) *model.User {
	a := (user.Email).String()

	return &model.User{
		Base:      model.Base{},
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     &a,
		Password:  *user.Password,
		IsAdmin:   false,
		Roles:     []*model.Role{},
	}
}

// LoginToUser converts a LoginUser to a User
func LoginToUser(user *api.LoginUser) *model.User {
	a := (user.Email).String()

	return &model.User{
		Email:    &a,
		Password: *user.Password,
	}
}
