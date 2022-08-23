package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Url  string `env:"CONNECTION_URL"`
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port string `env:"PORT" env-default:"8080"`
}

// Load config from enviroment
// Throw an error if broker connection string is not setted
func LoadConfig() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
