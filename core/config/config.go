package config

import (
	"log"
	"os"
	"path"
)

var Port string
var JwtSecret string
var DataDirectoryPath string
var DatabasePath string
var CorsBaseUrl string

func LoadConfig() {
	Port = getEnv("PORT", "8080")
	JwtSecret = getEnv("JWT_SECRET", "your-secret-key-change-this-in-production")
	DataDirectoryPath = getEnv("DATA_DIR_PATH", "./data")
	DatabasePath = getEnv("DB_PATH", path.Join(DataDirectoryPath, "dsn.db"))
	CorsBaseUrl = getEnv("CORS_BASE_URL", "*")

	if JwtSecret == "your-secret-key-change-this-in-production" {
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
