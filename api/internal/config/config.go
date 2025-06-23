package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("No .env file found, loading from environment variables. %w", err)
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, envErrorMsg("DB_USER")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, envErrorMsg("DB_PASSWORd")
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, envErrorMsg("DB_HOST")
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		return nil, envErrorMsg("DB_PORT")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, envErrorMsg("DB_NAME")
	}
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	port := os.Getenv("PORT")
	if port == "" {
		return nil, envErrorMsg("PORT")
	}

	return &Config{
		DatabaseURL: dbUrl,
		Port:        port,
	}, nil
}

func envErrorMsg(envStr string) error {
	if envStr == "" {
		envStr = "ENV_VARIABLE"
	}
	return fmt.Errorf("%s not set in environment", envStr)
}
