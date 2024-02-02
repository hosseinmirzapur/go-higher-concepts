package jwt

import "os"

type Config struct {
	SecretKey string
}

func JWTConfig() *Config {
	return &Config{
		SecretKey: os.Getenv("JWT_SECRET"),
	}
}
