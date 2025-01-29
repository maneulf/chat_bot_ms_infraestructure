package configs

import "github.com/kelseyhightower/envconfig"

type serviceConfig struct {
	DataBaseUrl string `envconfig:"DATA_BASE_URL"`
}

type Config struct {
	Service serviceConfig
}

func ConfigFromEnv(prefix string) *Config {
	APP := serviceConfig{}

	_ = envconfig.Process(prefix, &APP)
	return &Config{
		Service: APP,
	}
}
