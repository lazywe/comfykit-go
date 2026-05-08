package utils

import (
	"fmt"
	"os"
)

func GetEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func GetEnvIntOrDefault(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var val int
		if _, err := fmt.Sscanf(v, "%d", &val); err == nil {
			return val
		}
	}
	return def
}
