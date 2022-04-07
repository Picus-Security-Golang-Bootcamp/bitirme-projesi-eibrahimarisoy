package main

import (
	"fmt"
	"log"
	"net/http"
	"patika-ecommerce/pkg/config"
	db "patika-ecommerce/pkg/database"
	"patika-ecommerce/pkg/graceful"
	logger "patika-ecommerce/pkg/logging"
	"patika-ecommerce/pkg/router"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting server...")

	// Load the configuration file.

	cfg, err := config.LoadConfig("./pkg/config/config-local")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set logger
	logger.NewLogger(cfg)
	defer logger.Close()

	// Connect to database
	DB := db.NewPsqlDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB", DB)

	r := gin.Default()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ServerConfig.ReadTimeoutSecs) * time.Second,
		WriteTimeout: time.Duration(cfg.ServerConfig.WriteTimeoutSecs) * time.Second,
	}

	rootRouter := r.Group(cfg.ServerConfig.RoutePrefix)
	router.InitializeRoutes(rootRouter, DB, cfg)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()
	log.Println("Patika ecommerce service started")

	// Wait for interrupt signal to gracefully shutdown the server with
	graceful.ShutdownGin(srv, time.Duration(cfg.ServerConfig.TimeoutSecs*int64(time.Second)))
}
