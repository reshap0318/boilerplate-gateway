package helpers

import (
	"os"
	"strconv"
)

// GetEnv retrieves an environment variable or returns a default value if not set.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retrieves an environment variable as integer or returns a default value.
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// MessageBaseURL is where serv-message is reachable.
func MessageBaseURL() string {
	return GetEnv("SVC_MESSAGE_URL", "http://message:8080")
}
