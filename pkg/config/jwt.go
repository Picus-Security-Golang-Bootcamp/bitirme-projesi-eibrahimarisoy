package config

type JWTConfig struct {
	SecretKey           string
	AccessTokenLifeTime int
	RefrehTokenLifeTime int
}
