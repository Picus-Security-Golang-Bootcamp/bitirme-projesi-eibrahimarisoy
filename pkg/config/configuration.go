package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig ServerConfig
	JWTConfig    JWTConfig
	DBConfig     DatabaseConfig
}

// LoadConfig loads the configuration from the given file.
func LoadConfig(file string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(file)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}

		return nil, err
	}

	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
