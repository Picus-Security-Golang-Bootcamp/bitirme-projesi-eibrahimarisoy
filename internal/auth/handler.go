package auth

import (
	"fmt"
	"patika-ecommerce/internal/api"
	httpErr "patika-ecommerce/internal/httpErrors"
	"patika-ecommerce/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type authHandler struct {
	authService AuthServiceInterface
	cfg         *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(r *gin.RouterGroup, cfg *config.Config, authService AuthServiceInterface) {
	handler := &authHandler{
		cfg:         cfg,
		authService: authService,
	}

	r.POST("/register", handler.register)
	r.POST("/login", handler.login)
	r.POST("/refresh", handler.refreshToken)
}

// register is used to register a new user
func (u *authHandler) register(c *gin.Context) {
	var reqBody api.RegisterUser
	fmt.Println("Register121212121121")

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	resp, err := u.authService.Register(RegisterToUser(&reqBody))

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(201, resp)
}

// login is used to login a user
func (u *authHandler) login(c *gin.Context) {
	var reqBody api.LoginUser

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	resp, err := u.authService.Login(LoginToUser(&reqBody))

	if err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	c.JSON(200, resp)
}

// refresh is used to refresh the token
func (u *authHandler) refreshToken(c *gin.Context) {
	var reqBody api.RefreshToken
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	if err := reqBody.Validate(strfmt.NewFormats()); err != nil {
		fmt.Println(err)
		c.JSON(httpErr.ErrorResponse(err))
		return
	}

	resp, err := u.authService.RefreshToken(*reqBody.RefreshToken)
	fmt.Println(err)
	if err != nil {
		fmt.Println(err)

		c.JSON(httpErr.ErrorResponse(err))
		return
	}
	fmt.Println("dasdasdasodk446546465465465")

	c.JSON(200, resp)
}
