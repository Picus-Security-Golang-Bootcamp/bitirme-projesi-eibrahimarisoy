package jwt_helper

import (
	"encoding/json"

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

func GenerateToken(claims *jwt.Token, secret string) string {
	hmacSecretString := secret
	hmacSecret := []byte(hmacSecretString)
	token, _ := claims.SignedString(hmacSecret)

	return token
}

func VerifyToken(token string, secret string) *JWTToken {
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

	return &decodedToken
}
