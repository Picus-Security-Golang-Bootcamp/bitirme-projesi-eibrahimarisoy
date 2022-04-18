package auth

import (
	"errors"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	user "patika-ecommerce/internal/user"
	"patika-ecommerce/pkg/config"
	jwtHelper "patika-ecommerce/pkg/jwt"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30},
	}

	type fields struct {
		cfg      *config.Config
		userRepo user.UserRepositoryInterface
	}
	type args struct {
		user *model.User
	}

	firstname, lastname, email, username, password := "test", "test", "test@example.com", "test_test", "123456Aa"
	wrongEmail := "wrong@email.com"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockRepo := &mockUserRepository{
		items: []model.User{
			{
				Base:      model.Base{ID: uuid.New()},
				FirstName: &firstname,
				LastName:  &lastname,
				Email:     &email,
				Password:  string(hashed),
			},
		},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.TokenResponse
		wantErr bool
	}{
		{
			name: "userLogin_Success",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname,
					LastName:  &lastname,
					Username:  &username,
					Email:     &email,
					Password:  password,
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: false,
		},
		{
			name: "userLogin_Failed_wrongPassword",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname,
					LastName:  &lastname,
					Username:  &username,
					Email:     &email,
					Password:  "123456",
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "userLogin_Failed_wrongEmail",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname,
					LastName:  &lastname,
					Username:  &username,
					Email:     &wrongEmail,
					Password:  password,
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				cfg:      tt.fields.cfg,
				userRepo: tt.fields.userRepo,
			}
			_, err := a.Login(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30},
	}
	firstname, lastname, email, username, password := "test", "test", "test@example.com", "test_test", "123456Aa"
	firstname2, lastname2, email2, username2, password2 := "test2", "test2", "test2@example.com", "test2_test2", "123456Aa"

	mockRepo := &mockUserRepository{
		items: []model.User{
			{
				Base:      model.Base{ID: uuid.New()},
				FirstName: &firstname,
				LastName:  &lastname,
				Username:  &username,
				Email:     &email,
				Password:  password,
			},
		},
	}

	type fields struct {
		cfg      *config.Config
		userRepo user.UserRepositoryInterface
	}
	type args struct {
		user *model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.TokenResponse
		wantErr bool
	}{
		{
			name: "userRegister_Successful",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname2,
					LastName:  &lastname2,
					Username:  &username2,
					Email:     &email2,
					Password:  password2,
					IsAdmin:   false,
				},
			},
			want:    api.TokenResponse{},
			wantErr: false,
		},
		{
			name: "userRegister_Failed_duplicateEmail",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname2,
					LastName:  &lastname2,
					Username:  &username2,
					Email:     &email,
					Password:  "123456",
					IsAdmin:   true,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "userRegister_Failed_duplicateUsername",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				user: &model.User{
					Base:      model.Base{ID: uuid.New()},
					FirstName: &firstname2,
					LastName:  &lastname2,
					Username:  &username,
					Email:     &email2,
					Password:  "123456",
					IsAdmin:   true,
				},
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				cfg:      tt.fields.cfg,
				userRepo: tt.fields.userRepo,
			}
			_, err := a.Register(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30},
	}
	firstname, lastname, email, username, password := "test", "test", "test@example.com", "test", "123456Aa"

	mockRepo := &mockUserRepository{
		items: []model.User{
			{
				Base:      model.Base{ID: uuid.New()},
				FirstName: &firstname,
				LastName:  &lastname,
				Username:  &username,
				Email:     &email,
				Password:  password,
			},
		},
	}
	invalidRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	expiredRefreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDcmVhdGVkQXQiOjE2NTAwMjUyMTYsIkV4cGlyZXNBdCI6MTY1MDAyODgxNiwiVXNlcklkIjoiNWQ0ODBmZTItYzY2Yy00MzkwLWJhMTctNmJmZTMxMTc5ZDY5In0.HyfE06jGkTNPyBXFMQwyUa1lGs83xsSv5N7my4kWtYY"
	anonymousUserToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJDcmVhdGVkQXQiOjE2NTAwMjUyMTYsIkV4cGlyZXNBdCI6MTY1MDAyODgxNiwiVXNlcklkIjoiNWQ0ODBmZTItYzY2Yy00MzkwLWJhMTctNmJmZTMxMTc5ZDkifQ.OWIWseL90E9h3ainPmcqj6cgUvWZDoUQZ4xEzUbxTPk"

	type fields struct {
		cfg      *config.Config
		userRepo user.UserRepositoryInterface
	}
	type args struct {
		refreshToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    api.TokenResponse
		wantErr bool
	}{
		{
			name: "refreshToken_Successful",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: jwtHelper.GetAuthToken(&mockRepo.items[0], cfg).RefreshToken,
			},
			want:    api.TokenResponse{},
			wantErr: false,
		},
		{
			name: "refreshToken_Failed_InvalidRefreshToken",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: invalidRefreshToken,
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "refreshToken_Failed_ExpiredRefreshToken",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: expiredRefreshToken,
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
		{
			name: "refreshToken_Failed_AnonymousUser",
			fields: fields{
				cfg:      cfg,
				userRepo: mockRepo,
			},
			args: args{
				refreshToken: anonymousUserToken,
			},
			want:    api.TokenResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthService{
				cfg:      tt.fields.cfg,
				userRepo: tt.fields.userRepo,
			}
			_, err := a.RefreshToken(tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// Mock UserRepository
var (
	errCRUD      = errors.New("Mock: Error crud operation")
	userNotFound = errors.New("User not found")
)

type mockUserRepository struct {
	items []model.User
}

// InsertUser insert user to mock repository
func (u *mockUserRepository) InsertUser(user *model.User) (*model.User, error) {
	for _, item := range u.items {
		if *item.Username == *user.Username || *item.Email == *user.Email {
			return nil, errCRUD
		}
	}
	u.items = append(u.items, *user)
	return user, nil
}

// GetUser get user from mock repository
func (u *mockUserRepository) GetUser(id string) (*model.User, error) {

	if len(id) == 0 {
		return nil, errCRUD
	}
	uu, _ := uuid.Parse(id)

	for _, item := range u.items {
		if item.ID == uu {
			return &item, nil
		}
	}
	return nil, userNotFound
}

// GetUserByEmail get user by email from mock repository
func (u *mockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	if len(email) == 0 {
		return nil, errCRUD
	}
	for _, item := range u.items {
		if *item.Email == email {
			return &item, nil
		}
	}
	return &user, nil
}
