package auth

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"
	jwtHelper "patika-ecommerce/pkg/jwt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtToken struct {
	accesToken   string
	refreshToken string
}

type AuthService struct {
	cfg *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		cfg: cfg,
	}
}

//AuthTokenService is a service that generates a new JWT token
func (a *AuthService) GetAuthTokenService(user model.User) api.TokenResponse {
	jwtClaimsForAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"CreatedAt": time.Now().Unix(),
		"ExpiresAt": time.Now().Add(time.Duration(a.cfg.JWTConfig.AccessTokenLifeTime) * time.Hour).Unix(),
		"Roles":     user.Roles,
		"UserId":    user.ID,
		"Email":     user.Email,
		"IsAdmin":   user.IsAdmin,
	})

	jwtClaimsForRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"CreatedAt": time.Now().Unix(),
		"ExpiresAt": time.Now().Add(time.Duration(a.cfg.JWTConfig.RefrehTokenLifeTime) * time.Hour).Unix(),
		"UserId":    user.ID,
	})

	accesToken := jwtHelper.GenerateToken(jwtClaimsForAccessToken, a.cfg.JWTConfig.SecretKey)
	refreshToken := jwtHelper.GenerateToken(jwtClaimsForRefreshToken, a.cfg.JWTConfig.SecretKey)

	return api.TokenResponse{
		AccessToken:  accesToken,
		RefreshToken: refreshToken,
	}
}
