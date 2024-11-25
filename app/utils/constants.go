package utils

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var SERVER_PORT = getEnvInt("LISTEN_PORT", 3000)

var DB_HOST = getEnvStr("DB_HOSTNAME", "")
var DB_PORT = getEnvInt("DB_PORT", 0000)
var DB_USER = getEnvStr("DB_USERNAME", "")
var DB_PASSWORD = getEnvStr("DB_PASSWORD", "")
var DB_NAME = getEnvStr("DB_NAME", "")

var SECRET_KEY = getEnvStr("SECRET_KEY", "")

// Get an Environment Variable as an integer value
func getEnvInt(envVar string, defaultVal int) int {
	intEnvVar, exists := os.LookupEnv(envVar)
	if !exists {
		fmt.Println("[utils]", envVar, "does not exist. Will use default value:", defaultVal)
		return defaultVal
	}

	envInteger, err := strconv.Atoi(intEnvVar)
	if err != nil {
		fmt.Println("[utils]", err.Error(), "Will use default value:", defaultVal)
		return defaultVal
	}

	return envInteger
}

// Get an Environment Variable as a string value
func getEnvStr(envVar string, defaultVal string) string {
	strEnvVar, exists := os.LookupEnv(envVar)
	if !exists {
		fmt.Println("[utils]", envVar, "does not exist. Will use default value:", defaultVal)
		return defaultVal
	}

	return strEnvVar
}
