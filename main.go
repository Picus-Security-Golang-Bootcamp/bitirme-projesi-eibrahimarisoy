package main

import (
	"log"
	"patika-ecommerce/pkg/config"
	db "patika-ecommerce/pkg/database"
	logger "patika-ecommerce/pkg/logging"
)

func main() {
	log.Println("Starting server...")

	// Load the configuration file.

	cfg, err := config.LoadConfig("config-local.yaml")
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
}
