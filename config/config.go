package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

// Config holds all the environment variables for the application.
type Config struct {
	DBHost string `validate:"required"`
	DBPort string `validate:"required"`
	DBUser string `validate:"required"`
	DBPassword string `validate:"required"`
	DBName string `validate:"required"`
	SSLMode string `validate:"required"`
	ServerPort string `validate:"required"`
	TokenTTL string `validate:"required"`
	AccessTokenSecret string `validate:"required"`
}

// LoadENV loads configuration from .env file and environment variables.
func LoadENV() (*Config, error) {
	
	// Load the .env file if it exists
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	/*
	// Load .env file for the current environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development" // Default to development if not set
	}
	err := godotenv.Load(fmt.Sprintf(".env.%s", env))
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file for %s environment: %v", env, err)
	}
	*/

	// Load environment variables into a Config struct
	config := &Config{
		DBHost:		os.Getenv("DB_HOST"),
		DBPort: 	os.Getenv("DB_PORT"),
		DBUser: 	os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: 	os.Getenv("DB_NAME"),
		SSLMode:	os.Getenv("SSL_MODE"),
		ServerPort: os.Getenv("SERVER_PORT"),
		TokenTTL:	os.Getenv("TOKEN_TTL"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),

	}

	// Validate configuration
	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("Configuration validation failed: %v", err)
	}

	return config, nil
}

