package config

type Config struct {
	GRPCAddress string
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	GRPCAddress := getenv("GRPC_ADDRESS")

	return Config{
		GRPCAddress,
	}, nil
}
