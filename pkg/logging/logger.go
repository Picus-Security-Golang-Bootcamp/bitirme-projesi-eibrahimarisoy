package logger

import (
	"fmt"
	"patika-ecommerce/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(config *config.Config) {
	logLevel, err := zapcore.ParseLevel(config.LoggerConfig.Level)

	if err != nil {
		panic(fmt.Sprintf("Unknown log level: %v", logLevel))
	}

	var cfg zap.Config

	if config.LoggerConfig.Development {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	logger, err := cfg.Build()

	if err != nil {
		logger = zap.NewNop()
	}
	zap.ReplaceGlobals(logger)

}

func Close() {
	zap.L().Sync()
}
