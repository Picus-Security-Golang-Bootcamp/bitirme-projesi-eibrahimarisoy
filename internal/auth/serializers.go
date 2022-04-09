package auth

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
)

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

func LoginToUser(user *api.LoginUser) *model.User {
	a := (user.Email).String()

	return &model.User{
		Email:    &a,
		Password: *user.Password,
	}
}
