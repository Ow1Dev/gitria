package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Config struct {
	Data      string
	Ssh  			SSHConfig
	HTTP 			HTTPConfig
}

func (c *Config) GetDbFilePath() string {
	return filepath.Join(c.Data, "store.db") 
}

func (c *Config) GetRepoFolderPath() string {
	return filepath.Join(c.Data, "repositories")  
}

type SSHConfig struct {
	HostKeyPath string
	PORT 				string
}

type HTTPConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	datapath := getenv("DATA")
	if datapath == "" {
		datapath = "./data"
	}

	addr := getAddr()
	return Config{
		Data:  datapath,
		Ssh: SSHConfig{
			PORT: 				 getenv("SSH_LISTEN_PORT"),
			HostKeyPath:   getenv("SSH_HOST_KEY"),
		},
		HTTP: HTTPConfig{
			Address:      addr,
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		},
	}, nil
}

func getAddr() string {
	if a := os.Getenv("ADDRESS"); a != "" {
		return a
	}
	if p := os.Getenv("PORT"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			return ":" + p
		}
		return p
	}
	return ":8080"
}

func getDurationEnv(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
