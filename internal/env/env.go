package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	env := os.Getenv(key)
	if env == "" {
		return fallback
	}
	return env
}

func GetInt(key string, fallback int) (int, error) {
	env := os.Getenv(key)
	if env == "" {
		return fallback, nil
	}
	valAsInt, err := strconv.Atoi(env)
	if err != nil {
		return fallback, err
	}
	return valAsInt, nil
}
