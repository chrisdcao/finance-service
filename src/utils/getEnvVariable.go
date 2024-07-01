package utils

import (
	"fmt"
	"os"
)

// GetEnv fetches an environment variable and validates it.
func GetEnvVariableFromKey(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is not set", key)
	}
	return value, nil
}
