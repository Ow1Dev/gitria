package config

type Config struct {
	Ssh SSHConfig
}

type SSHConfig struct {
	HostKeyPath string
	PORT string
}

func LoadFromEnv(getenv func(string) string) (Config, error) {
	return Config{
		SSHConfig{
			PORT: 				 getenv("SSH_LISTEN_PORT"),
			HostKeyPath:   getenv("SSH_HOST_KEY"),
		},
	}, nil
}
