package configs

import "github.com/kelseyhightower/envconfig"

type serviceConfig struct {
	CsmlPath               string `envconfig:"CSML_PATH"`
	CsmlXApiKey            string `envconfig:"CSML_X_API_KEY"`
	DataBaseMessageSaverMS string `envconfig:"DATA_BASE_MESSAGES_SAVER_MS"`
}

type Config struct {
	Service serviceConfig
}

func ConfigFromEnv(prefix string) *Config {
	APP := serviceConfig{
		CsmlPath: "",
	}

	_ = envconfig.Process(prefix, &APP)
	return &Config{
		Service: APP,
	}
}
