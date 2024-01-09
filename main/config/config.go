package config

import "os"

const (
	prod = "production"
)

type Config struct {
	Env       string `env:"ENV"`
	Host      string `env:"APP_HOST"`
	Port      string `env:"APP_PORT"`
	JWTSecret string `env:"JWT_SIGN_KEY"`
	HMACKey   string `env:"HMAC_KEY"`
}

func (c Config) IsProd() bool {
	return c.Env == prod
}

func GetConfig() Config {
	return Config{
		Env:       os.Getenv("ENV"),
		Host:      os.Getenv("APP_HOST"),
		Port:      os.Getenv("APP_PORT"),
		JWTSecret: os.Getenv("JWT_SIGN_KEY"),
		HMACKey:   os.Getenv("HMAC_KEY"),
	}
}
