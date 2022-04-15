package auth

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"patika-ecommerce/pkg/config"
	jwtHelper "patika-ecommerce/pkg/jwt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func getRegisterPOSTPayload() []byte {
	var jsonStr = []byte(
		`{"firstName":"emre","lastName":"yilmaz","email":"e@emre.com", "password":"123456", "username":"emre1", "isAdmin":"false"}`)

	return jsonStr
}

func Test_authHandler_register(t *testing.T) {
	firstname, lastname, email, username, password := "test", "test", "emre@arisoy.com", "emre", "123456Aa"

	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30},
	}
	mockAuthService := &mockAuthService{
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
		cfg: cfg,
	}
	w := httptest.NewRecorder()
	authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
	c, r := gin.CreateTestContext(w)

	r.POST("/register", authHandler.register)
	c.Request, _ = http.NewRequest("POST", "/register", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getRegisterPOSTPayload()))
	authHandler.register(c)

	assert.Equal(t, http.StatusCreated, w.Code)

}

func Test_authHandler_login(t *testing.T) {
	firstname, lastname, email, username, password := "test", "test", "emre@arisoy.com", "emre", "123456Aa"
	loginPayload := []byte(`{"email":"emre@arisoy.com","password": "123456Aa"}`)

	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30},
	}
	mockAuthService := &mockAuthService{
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
		cfg: cfg,
	}
	w := httptest.NewRecorder()
	authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
	c, r := gin.CreateTestContext(w)

	r.POST("/login", authHandler.login)
	c.Request, _ = http.NewRequest("POST", "/login", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(loginPayload))
	authHandler.login(c)

	assert.Equal(t, http.StatusOK, w.Code)

}

func Test_authHandler_refreshToken(t *testing.T) {
	firstname, lastname, email, username, password := "test", "test", "emre@arisoy.com", "emre", "123456Aa"

	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30, RefreshTokenLifeTime: 30},
	}
	mockAuthService := &mockAuthService{
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
		cfg: cfg,
	}
	token := jwtHelper.GetAuthToken(&mockAuthService.items[0], cfg).RefreshToken
	refreshTokenPayload := []byte(`{"refreshToken":"` + token + `"}`)

	w := httptest.NewRecorder()
	authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
	c, r := gin.CreateTestContext(w)

	r.POST("/refresh", authHandler.refreshToken)
	c.Request, _ = http.NewRequest("POST", "/refresh", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(refreshTokenPayload))
	authHandler.refreshToken(c)
	fmt.Println(w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

}

type mockAuthService struct {
	items []model.User
	cfg   *config.Config
}

var UsernameAlreadyExists = fmt.Errorf("23505")

// Register is a service that registers a new user
func (a *mockAuthService) Register(user *model.User) (api.TokenResponse, error) {
	for _, item := range a.items {
		if *item.Username == *user.Username || *item.Email == *user.Email {
			return api.TokenResponse{}, UsernameAlreadyExists
		}
	}
	a.items = append(a.items, *user)

	return jwtHelper.GetAuthToken(user, a.cfg), nil
}

// Login is a service that logs in a user
func (a *mockAuthService) Login(u *model.User) (api.TokenResponse, error) {

	for _, item := range a.items {
		if *item.Email == *u.Email {
			if item.Password == u.Password {
				return jwtHelper.GetAuthToken(&item, a.cfg), nil
			}
		}
	}
	return api.TokenResponse{}, fmt.Errorf("Invalid username or password")
}

// RefreshToken is a service that checks if the refresh token is valid and returns a new JWT token
func (a *mockAuthService) RefreshToken(refreshToken string) (api.TokenResponse, error) {
	fmt.Println("rewrewrwerwe")

	token, err := jwtHelper.ParseToken(refreshToken, a.cfg.JWTConfig.SecretKey)

	if err != nil {
		fmt.Println("rewrewrwerwe")

		return api.TokenResponse{}, err
	}

	claims := token.Claims.(jwt.MapClaims)

	if int64(claims["ExpiresAt"].(float64)) < time.Now().Unix() {
		return api.TokenResponse{}, fmt.Errorf("refresh token expired")
	}

	userId := claims["UserId"].(string)

	id, _ := uuid.Parse(userId)
	user := &model.User{}
	for _, item := range a.items {
		if item.ID == id {
			user = &item
			break
		}
	}

	if user == nil {
		fmt.Println("dddddddddddddd")
		return api.TokenResponse{}, err
	}

	jwtClaimsForAccessToken := jwtHelper.NewJwtClaimsForAccessToken(user, a.cfg.JWTConfig.AccessTokenLifeTime)

	accesToken := jwtHelper.GenerateToken(jwtClaimsForAccessToken, a.cfg.JWTConfig.SecretKey)
	fmt.Println("fsdfsdfssssseddddddddddddddddddddddd")

	return api.TokenResponse{
		AccessToken: accesToken,
	}, nil
}
