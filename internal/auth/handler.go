package auth

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/pkg/config"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService *AuthService
	cfg         *config.Config
}

func NewAuthHandler(r *gin.RouterGroup, authService *AuthService, cfg *config.Config) {
	handler := &authHandler{
		authService: authService,
		cfg:         cfg,
	}

	r.POST("/register", handler.register)
	r.POST("/login", handler.login)
	r.POST("/refresh", handler.refresh)
}

func (u *authHandler) register(c *gin.Context) {
	var userBody api.RegisterUser
	if err := c.ShouldBindJSON(&userBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := u.authService.RegisterService(RegisterToUser(&userBody))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (u *authHandler) login(c *gin.Context) {
	var reqBody api.LoginUser

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	resp, err := u.authService.LoginService(LoginToUser(&reqBody))

	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(200, resp)
}

// refresh is used to refresh the token
func (u *authHandler) refresh(c *gin.Context) {
	var reqBody api.RefreshToken
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	resp, err := u.authService.RefreshTokenService(*reqBody.RefreshToken)
	// user, err := u.userRepo.GetUserByEmail((reqBody.Email).String())
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(200, resp)
}
