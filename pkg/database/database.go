package db

import (
	"fmt"
	"patika-ecommerce/pkg/config"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPsqlDB creates a new database connection
func NewPsqlDB(cfg *config.Config) *gorm.DB {

	fmt.Println("dataSourceName: ", cfg.DBConfig.DataSourceName)
	db, err := gorm.Open(postgres.Open(cfg.DBConfig.DataSourceName), &gorm.Config{})

	if err != nil {
		zap.L().Fatal("Error connecting to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Fatal("Error connecting to database", zap.Error(err))
	}

	if err := sqlDB.Ping(); err != nil {
		zap.L().Fatal("Cannot ping database", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(cfg.DBConfig.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.DBConfig.MaxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConfig.MaxLifetime) * time.Second)

	return db
}
