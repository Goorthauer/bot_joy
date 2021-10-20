package config

import (
	"os"
	"strconv"
	"strings"
)

type QueryConfig struct {
	Call     []string
	Response []string
}

type Config struct {
	Query []QueryConfig
	Token string
}

// New returns a new Config struct
func New() *Config {
	allCall := Explode(",", getEnv("ALL_QUERY", ""))
	allResponse := Explode(",", getEnv("ALL_RESPONSE", ""))
	var newQueryConfig = make([]QueryConfig, len(allCall))
	for i, callQuery := range allCall {
		responseQuery := allResponse[i]
		newQueryConfig[i] = QueryConfig{
			Call:     Explode(",", getEnv(callQuery, "")),
			Response: Explode(",", getEnv(responseQuery, "")),
		}
	}
	return &Config{
		Query: newQueryConfig,
		Token: getEnv("TOKEN", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}
