package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddress string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	RepoAddr string
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	HTTPPort := getenv("PORT")
	if HTTPPort == "" {
		HTTPPort = "9000"
	}

	repoAddr := getenv("GIT_REPO_GRPC")
	if repoAddr == "" {
		return Config{}, fmt.Errorf("GIT_REPO_GRPC is required")
	}

	return Config{
		HTTPAddress: fmt.Sprintf(":%s", HTTPPort),
		ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		RepoAddr: repoAddr,
	}, nil
}

func getDurationEnv(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
