package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerAddress  string
	PollInterval   int
	ReportInterval int
}

func LoadConfig() Config {
	address := getenv("ADDRESS", "localhost:8080")
	pollInterval := getenvInt("POLL_INTERVAL", 2)
	reportInterval := getenvInt("REPORT_INTERVAL", 10)

	return Config{
		ServerAddress:  address,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
