package jwt_helper

import (
	"encoding/json"
	"patika-ecommerce/internal/model"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/golang-jwt/jwt/v4"
)

type JWTToken struct {
	CreatedAt int      `json:"created_at"`
	ExpiresAt int      `json:"expires_at"`
	Roles     []string `json:"roles"`
	UserId    string   `json:"userId"`
	Email     string   `json:"email"`
	IsAdmin   bool     `json:"isAdmin"`
}

// GenerateAccessToken generates an access token
func GenerateToken(claims *jwt.Token, secret string) string {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)
	token, _ := claims.SignedString(hmacSecret)

	return token
}

//VerifyToken verifies a token
func VerifyToken(token string, secret string) *model.User {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)

	decoded, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil
	}

	if !decoded.Valid {
		return nil
	}

	decodedClaims := decoded.Claims.(jwt.MapClaims)

	var decodedToken JWTToken
	jsonString, _ := json.Marshal(decodedClaims)
	json.Unmarshal(jsonString, &decodedToken)

	// decode the claims into a model.user struct and return it
	user := model.User{
		Base: model.Base{
			ID: strfmt.UUID(decodedToken.UserId),
		},
		Email:   &decodedToken.Email,
		IsAdmin: decodedToken.IsAdmin,
		Roles:   []*model.Role{},
	}

	return &user
}

// ParseToken parses a token string into a jwt.Token
func ParseToken(token string, secret string) (*jwt.Token, error) {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
}

// NewJwtClaimsForAccessToken returns a new jwt.Claims for an access token
func NewJwtClaimsForAccessToken(user *model.User, accessTokenLifeTime int) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"CreatedAt": time.Now().Unix(),
		"ExpiresAt": time.Now().Add(time.Duration(accessTokenLifeTime) * time.Hour).Unix(),
		"UserId":    user.ID,
		"Email":     user.Email,
		"IsAdmin":   user.IsAdmin,
		"Roles":     user.Roles,
	})
}

// NewJwtClaimsForRefreshToken returns a new jwt.Claims for a refresh token
func NewJwtClaimsForRefreshToken(user *model.User, refreshTokenLifeTime int) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"CreatedAt": time.Now().Unix(),
		"ExpiresAt": time.Now().Add(time.Duration(refreshTokenLifeTime) * time.Hour).Unix(),
		"UserId":    user.ID,
	})
}
