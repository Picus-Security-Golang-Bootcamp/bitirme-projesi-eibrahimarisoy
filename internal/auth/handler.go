package auth

import (
	"fmt"
	"patika-ecommerce/internal/api"
	user "patika-ecommerce/internal/user"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	userRepo    *user.UserRepository
	authService *AuthService
}

func NewAuthHandler(r *gin.RouterGroup, userRepo *user.UserRepository, authService *AuthService) {
	handler := &authHandler{
		userRepo:    userRepo,
		authService: authService,
	}

	r.POST("/register", handler.register)
	r.POST("/login", handler.login)
}

func (u *authHandler) register(c *gin.Context) {
	var userBody api.RegisterUser
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := u.userRepo.InsertUser(RegisterToUser(&userBody))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, u.authService.GetAuthTokenService(*user))
}

func (u *authHandler) login(c *gin.Context) {
	var reqBody api.LoginUser

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	user, err := u.userRepo.GetUserByEmail((reqBody.Email).String())
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	// user := LoginToUser(&reqBody)

	if !user.CheckPassword(*reqBody.Password) {
		c.JSON(400, gin.H{"msg": "Invalid password"})
		return
	}
	fmt.Println(LoginToUser(&reqBody))

	// if err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(200, u.authService.GetAuthTokenService(*user))
}
