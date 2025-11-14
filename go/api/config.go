package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	UsersTable         string

	JWTSecret string
}

func Load() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}

	cfg := &Config{
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),

		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		UsersTable:         os.Getenv("USERS_TABLE"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" || cfg.GoogleRedirectURL == "" {
		log.Fatal("Google OAuth env vars missing")
	}
	if cfg.AWSRegion == "" || cfg.UsersTable == "" {
		log.Fatal("AWS env vars missing")
	}

	return cfg
}
