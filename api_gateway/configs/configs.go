package configs

import "github.com/kelseyhightower/envconfig"

type serviceConfig struct {
	WebAdapterMSPath string `envconfig:"WEB_ADAPTER_MS_PATH"`
}

type Config struct {
	Service serviceConfig
}

func ConfigFromEnv(prefix string) *Config {
	APP := serviceConfig{
		WebAdapterMSPath: "",
	}

	_ = envconfig.Process(prefix, &APP)
	return &Config{
		Service: APP,
	}
}
