package config

import (
	"cmp"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/joho/godotenv"
)

var Port int
var JwtSecret string
var DataDirectoryPath string
var CorsBaseUrl string
var DatabaseDirectory string
var UploadsDirectory string
var NoAuthForUserZero bool

var defaults = map[string]string{
	"PORT":                  "8080",
	"DATA_DIR_PATH":         "./data",
	"CORS_BASE_URL":         "*",
	"AUTH_ENCRYPTION_KEY":   "0123456789abcdef0123456789abcdef",
	"NO_AUTH_FOR_USER_ZERO": "false",
}

func LoadConfig() {
	godotenv.Load(".env")
	var err error

	port := getEnv("PORT")
	Port, err = strconv.Atoi(port)
	if err != nil {
		log.Printf("*** Invalid PORT environment variable '%s', using default %s", port, defaults["PORT"])
		Port, _ = strconv.Atoi(defaults["PORT"])
	}

	DataDirectoryPath = getEnv("DATA_DIR_PATH")

	CorsBaseUrl = getEnv("CORS_BASE_URL")
	if CorsBaseUrl == defaults["CORS_BASE_URL"] {
		log.Println("*** CORS_BASE_URL environment variable is not set, allowing all origins")
	}

	JwtSecret = getEnv("AUTH_ENCRYPTION_KEY")
	if JwtSecret == defaults["AUTH_ENCRYPTION_KEY"] {
		log.Println("*** AUTH_ENCRYPTION_KEY environment variable is not set, using fallback key")
	}

	UploadsDirectory = path.Join(DataDirectoryPath, "uploads")
	DatabaseDirectory = path.Join(DataDirectoryPath, "database")

	NoAuthForUserZero = getEnv("NO_AUTH_FOR_USER_ZERO") == "true"
	if NoAuthForUserZero {
		log.Println("*** NO_AUTH_FOR_USER_ZERO is enabled, authentication middleware is disabled")
	}
}

func getEnv(environmentVariable string) string {
	return cmp.Or(os.Getenv(environmentVariable), defaults[environmentVariable])
}
