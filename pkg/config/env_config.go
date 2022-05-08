package config

import "os"

func GetEnvConfig(envKey string) string {
	envValue := os.Getenv(envKey)
	return envValue
}
