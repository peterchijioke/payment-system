package config

import (
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	WebhookSecret   string
	KorapayBaseURL  string
	RateLimitLimit  int
	RateLimitWindow time.Duration
}

func LoadConfig() *ServerConfig {

	return &ServerConfig{
		Port:            ":8080",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		WebhookSecret:   os.Getenv("WEBHOOK_SECRET"),
		KorapayBaseURL:  "https://checkout.korapay.com",
		RateLimitLimit:  getEnvInt("RATE_LIMIT", 300),
		RateLimitWindow: getEnvDuration("RATE_LIMIT_WINDOW", time.Minute),
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return time.Duration(intVal) * time.Second
		}
	}
	return defaultValue
}
