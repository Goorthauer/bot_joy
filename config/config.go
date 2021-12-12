package config

import (
	"os"
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

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	}
	return strings.Split(text, delimiter)

}
