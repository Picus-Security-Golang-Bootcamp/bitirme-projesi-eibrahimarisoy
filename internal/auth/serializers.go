package auth

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"

	"github.com/go-openapi/strfmt"
)

func UserToResponse(user *model.User) *api.UserToRegisterResponse {
	return &api.UserToRegisterResponse{
		ID:        user.ID,
		Email:     strfmt.Email(*user.Email),
		Username:  *user.Username,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		IsAdmin:   user.IsAdmin,
	}
}

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
