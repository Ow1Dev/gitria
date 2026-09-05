package config

import "fmt"

type Config struct {
	GRPCAddress string
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	GRPCPort := getenv("GRPC_PORT")
	if GRPCPort == "" {
		GRPCPort = "9091"
	}

	return Config{
		GRPCAddress: fmt.Sprintf(":%s", GRPCPort),
	}, nil
}
