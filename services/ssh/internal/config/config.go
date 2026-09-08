package config

type Config struct {
	ListenAddress string
	HostKeyPath 	string
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	return Config{
		ListenAddress: getenv("SSH_LISTEN_ADDRESS"),
		HostKeyPath:   getenv("SSH_HOST_KEY"),
	}, nil
}
