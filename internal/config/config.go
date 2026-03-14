package config

import (
	"os"
	"strings"
	"time"
)

type ResponseMode string

const (
	ResponseModeSuccess ResponseMode = "success"
	ResponseModeError   ResponseMode = "error"
)

type Config struct {
	HTTPPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	ResponseMode ResponseMode
}

func MustLoad() Config {
	cfg := Config{
		HTTPPort:     getEnv("HTTP_PORT", ":8080"),
		ReadTimeout:  getDurationEnv("HTTP_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getDurationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getDurationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		ResponseMode: ResponseMode(strings.ToLower(getEnv("RESPONSE_MODE", "success"))),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}

	return d
}
