package auth_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/pkg/config"

	route "patika-ecommerce/pkg/router"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_authHandler_register(t *testing.T) {
	cfg, err := config.LoadConfig("../../pkg/config/config-local")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	mockdb, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockdb.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockdb,
	}), &gorm.Config{})

	// Connect to database
	// DB := db.NewPsqlDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	gin.SetMode(gin.TestMode)
	rootRouter := r.Group(cfg.ServerConfig.RoutePrefix)

	route.InitializeRoutes(rootRouter, gormDB, cfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(getRegisterPOSTPayload()))
	r.ServeHTTP(w, req)
	fmt.Println(w.Body.String())
	fmt.Println(w.Code)
	assert.Equal(t, 201, w.Code)
	// assert.Equal(t, "404 page not found", w.Body.String())

}

func getRegisterPOSTPayload() []byte {
	var jsonStr = []byte(
		`{"firstName":"emre","lastName":"yilmaz","email":"e@emre.com", "password":"123456", "username":"emre", "isAdmin":"false"}`)

	return jsonStr
}
