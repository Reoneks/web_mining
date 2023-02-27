package config

import (
	"encoding/json"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppConfig
	PostgresConfig
}

type AppConfig struct {
	AppAddr string `envconfig:"APP_SERVER_ADDRESS"`
}

type PostgresConfig struct {
	DSN      string `required:"true" default:"host=localhost user=postgres password=postgres dbname=backend-db port=5432 sslmode=disable" envconfig:"DB_DSN"`
}

func (cfg Config) String() string {
	buf, _ := json.MarshalIndent(&cfg, "", "")
	return string(buf)
}

var (
	once   sync.Once
	config *Config
)

func Get() (*Config, error) {
	var err error
	once.Do(func() {
		var cfg Config
		// If you run it locally and through terminal please set up this in Load function (../.env)
		_ = godotenv.Load(".env")

		if err = envconfig.Process("", &cfg); err != nil {
			return
		}

		config = &cfg
	})

	return config, err
}
