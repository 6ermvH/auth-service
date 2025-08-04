package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecretKey string
	DBUser       string
	DBPassword   string
	DBName       string
	DBHost       string
	DBPort       string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		JWTSecretKey: getEnv("JWT_SECRET_KEY", ""),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "password"),
		DBName:       getEnv("DB_NAME", "auth"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback != "" {
		log.Printf("Using fallback for %s", key)
		return fallback
	}
	log.Fatalf("FATAL: Environment variable %s is not set.", key)
	return ""
}
