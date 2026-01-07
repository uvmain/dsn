package config

import (
	"cmp"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

var Port string
var JwtSecret string
var DataDirectoryPath string
var CorsBaseUrl string
var DatabaseDirectory string
var UploadsDirectory string

func LoadConfig() {
	godotenv.Load(".env")

	Port = getEnv("PORT", "8080")
	JwtSecret = getEnv("AUTH_ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef")
	DataDirectoryPath = getEnv("DATA_DIR_PATH", "./data")
	CorsBaseUrl = getEnv("CORS_BASE_URL", "*")

	if JwtSecret == "0123456789abcdef0123456789abcdef" {
		log.Println("*** AUTH_ENCRYPTION_KEY environment variable is not set, using fallback key")
	}

	UploadsDirectory = path.Join(DataDirectoryPath, "uploads")
	DatabaseDirectory = path.Join(DataDirectoryPath, "database")
}

func getEnv(environmentVariable, defaultValue string) string {
	return cmp.Or(os.Getenv(environmentVariable), defaultValue)
}
