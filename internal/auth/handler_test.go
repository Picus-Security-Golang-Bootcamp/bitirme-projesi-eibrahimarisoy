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

var (
	firstname = "firstname"
	lastname  = "lastname"
	email     = "email@test.com"
	username  = "username"
	password  = "123456Aa"
)

func getRegisterPOSTPayloadSuccess() []byte {
	var jsonStr = []byte(
		`{"firstName":"test","lastName":"test","email":"test@test.com", "password":"123456", "username":"username", "isAdmin":"false"}`)

	return jsonStr
}

func getRegisterPOSTPayloadFault() []byte {
	var jsonStr = []byte(
		`{"first_Name":"test","last_Name":"test","email":"test@test.com", "password":"123456", "Username":"test_test", "isAdmin":"false"}`)

	return jsonStr
}

func Test_authHandler_register(t *testing.T) {
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30},
	}

	t.Run("userRegister_Successfully", func(t *testing.T) {

		mockAuthService := &mockAuthService{
			items: []model.User{},
			cfg:   cfg,
		}
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/register", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getRegisterPOSTPayloadSuccess()))
		authHandler.register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("userRegister_Failed_requestBody", func(t *testing.T) {

		mockAuthService := &mockAuthService{
			items: []model.User{},
			cfg:   cfg,
		}
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/register", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getRegisterPOSTPayloadFault()))
		authHandler.register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("userRegister_Failed_duplicateUsername", func(t *testing.T) {

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
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/register", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getRegisterPOSTPayloadSuccess()))
		authHandler.register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func Test_authHandler_login(t *testing.T) {
	loginPayloadSuccess := []byte(`{"email":"email@test.com","password": "123456Aa"}`)
	loginPayloadFaultRequestBody := []byte(`{"email":"test@test.com","password_test": "123456Aa"}`)
	loginPayloadFaultUserNotFound := []byte(`{"email":"email_fault@test.com","password": "123456Aa"}`)

	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30},
	}
	t.Run("userLogin_Successful", func(t *testing.T) {
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
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/login", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(loginPayloadSuccess))
		authHandler.login(c)

		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("userLogin_Failed_requestBody", func(t *testing.T) {
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
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/login", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(loginPayloadFaultRequestBody))
		authHandler.login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

	t.Run("userLogin_Failed_userNotFound", func(t *testing.T) {
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
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/login", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(loginPayloadFaultUserNotFound))
		authHandler.login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

}

func Test_authHandler_refreshToken(t *testing.T) {
	firstname, lastname, email, username, password := "test", "test", "emre@arisoy.com", "emre", "123456Aa"
	cfg := &config.Config{
		JWTConfig: config.JWTConfig{SecretKey: "secret", AccessTokenLifeTime: 30, RefreshTokenLifeTime: 30},
	}

	t.Run("refreshToken_Success", func(t *testing.T) {
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

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/refresh", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(refreshTokenPayload))
		authHandler.refreshToken(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("refreshToken_Failed_json", func(t *testing.T) {
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
		refreshTokenPayload := []byte(`"refreshToken":"` + token + `"`)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/refresh", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(refreshTokenPayload))
		authHandler.refreshToken(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("refreshToken_Failed_reqBodyFormat", func(t *testing.T) {
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
		refreshTokenPayload := []byte(`{"refresh__Token":"` + token + `"}`)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/refresh", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(refreshTokenPayload))
		authHandler.refreshToken(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("refreshToken_Failed_userNotFound", func(t *testing.T) {
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
		// token := jwtHelper.GetAuthToken(&mockAuthService.items[0], cfg).RefreshToken
		refreshTokenPayload := []byte(`{"refreshToken":"` + "token" + `"}`)

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		authHandler := &authHandler{authService: mockAuthService, cfg: cfg}
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/refresh", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(refreshTokenPayload))
		authHandler.refreshToken(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

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
	token, err := jwtHelper.ParseToken(refreshToken, a.cfg.JWTConfig.SecretKey)

	if err != nil {
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
		return api.TokenResponse{}, err
	}

	jwtClaimsForAccessToken := jwtHelper.NewJwtClaimsForAccessToken(user, a.cfg.JWTConfig.AccessTokenLifeTime)
	accesToken := jwtHelper.GenerateToken(jwtClaimsForAccessToken, a.cfg.JWTConfig.SecretKey)

	return api.TokenResponse{
		AccessToken: accesToken,
	}, nil
}
