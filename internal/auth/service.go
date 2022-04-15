package auth

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	user "patika-ecommerce/internal/user"
	"patika-ecommerce/pkg/config"
	jwtHelper "patika-ecommerce/pkg/jwt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	cfg      *config.Config
	userRepo *user.UserRepository
}

type AuthServiceInterface interface {
	Register(user *model.User) (api.TokenResponse, error)
	Login(user *model.User) (api.TokenResponse, error)
	GetAuthToken(user *model.User) api.TokenResponse
	RefreshToken(refreshToken string) (api.TokenResponse, error)
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *config.Config, userRepo *user.UserRepository) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

// Register is a service that registers a new user
func (a *AuthService) Register(user *model.User) (api.TokenResponse, error) {
	user, err := a.userRepo.InsertUser(user)
	if err != nil {
		return api.TokenResponse{}, err
	}
	return a.GetAuthToken(user), nil
}

// Login is a service that logs in a user
func (a *AuthService) Login(user *model.User) (api.TokenResponse, error) {
	user, err := a.userRepo.GetUserByEmail(*user.Email)
	if err != nil {
		return api.TokenResponse{}, err
	}

	if !user.CheckPassword(user.Password) {
		return api.TokenResponse{}, fmt.Errorf("invalid password")
	}

	return a.GetAuthToken(user), nil
}

// GetAuthToken is a service that generates a new JWT token
func (a *AuthService) GetAuthToken(user *model.User) api.TokenResponse {
	jwtClaimsForAccessToken := jwtHelper.NewJwtClaimsForAccessToken(user, a.cfg.JWTConfig.AccessTokenLifeTime)

	jwtClaimsForRefreshToken := jwtHelper.NewJwtClaimsForRefreshToken(user, a.cfg.JWTConfig.RefreshTokenLifeTime)

	accesToken := jwtHelper.GenerateToken(jwtClaimsForAccessToken, a.cfg.JWTConfig.SecretKey)
	refreshToken := jwtHelper.GenerateToken(jwtClaimsForRefreshToken, a.cfg.JWTConfig.SecretKey)

	return api.TokenResponse{
		AccessToken:  accesToken,
		RefreshToken: refreshToken,
	}
}

// RefreshToken is a service that checks if the refresh token is valid and returns a new JWT token
func (a *AuthService) RefreshToken(refreshToken string) (api.TokenResponse, error) {
	token, err := jwtHelper.ParseToken(refreshToken, a.cfg.JWTConfig.SecretKey)

	if err != nil {
		return api.TokenResponse{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	if int64(claims["ExpiresAt"].(float64)) > time.Now().Unix() {
		return api.TokenResponse{}, fmt.Errorf("refresh token expired")
	}

	userId := claims["UserId"].(string)
	user, err := a.userRepo.GetUser(userId)
	if err != nil {
		return api.TokenResponse{}, err
	}

	jwtClaimsForAccessToken := jwtHelper.NewJwtClaimsForAccessToken(user, a.cfg.JWTConfig.AccessTokenLifeTime)

	accesToken := jwtHelper.GenerateToken(jwtClaimsForAccessToken, a.cfg.JWTConfig.SecretKey)

	return api.TokenResponse{
		AccessToken: accesToken,
	}, nil
}
