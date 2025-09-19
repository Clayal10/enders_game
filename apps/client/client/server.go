package client

type ServerConfig struct {
	Hostname string
	Port     string
}

func Start(cfg *ServerConfig) {
	if err := validateConfig(cfg); err != nil {
		return
	}
}

func validateConfig(cfg *ServerConfig) error {
	return nil
}
