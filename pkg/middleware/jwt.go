package mw

import (
	"net/http"
	jwtHelper "patika-ecommerce/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthenticationMiddleware is a middleware that checks for valid JWT tokens
func AuthenticationMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			authHeader := c.GetHeader("Authorization")
			parsedAuthHeader := strings.Split(authHeader, " ")

			if parsedAuthHeader[0] == "Bearer" && len(parsedAuthHeader) == 2 {
				token := parsedAuthHeader[1]
				decodedClaims := jwtHelper.VerifyToken(token, secretKey)

				if decodedClaims != nil {
					c.Set("user", decodedClaims)

					c.Next()
					return
				}
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
			c.Abort()
			return

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized!"})
			c.Abort()
		}
		c.Abort()
		return
	}
}

// AdminMiddleware is a middleware that checks for request user is admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		isAdmin := user.(*jwtHelper.JWTToken).IsAdmin
		if !isAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized!"})
			c.Abort()
			return
		}

		c.Next()
	}
}
