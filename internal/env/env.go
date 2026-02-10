package env

import "os"

func Get(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func GetOptional(key string) (string, bool) {
	value := os.Getenv(key)

	if value == "" {
		return "", false
	}

	return value, true
}
