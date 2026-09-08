package config

type Config struct {
	PORT string
	HostKeyPath 	string
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	return Config{
		PORT: 				 getenv("SSH_LISTEN_PORT"),
		HostKeyPath:   getenv("SSH_HOST_KEY"),
	}, nil
}
