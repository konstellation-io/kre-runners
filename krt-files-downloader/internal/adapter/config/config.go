package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	VersionID   string `envconfig:"KRT_VERSION_ID" required:"true"`
	KRTBasePath string `envconfig:"KRT_BASE_PATH" required:"true"`
	MongoDB     struct {
		URI    string `envconfig:"KRT_MONGO_URI" required:"true"`
		DBName string `envconfig:"KRT_MONGO_DB_NAME" required:"true"`
		Bucket string `envconfig:"KRT_MONGO_BUCKET" required:"true"`
	}
}

func NewConfig() (*Config, error) {
	c := &Config{}

	err := envconfig.Process("", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
