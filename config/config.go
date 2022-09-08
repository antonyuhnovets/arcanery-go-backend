package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer
	Database
}

type HttpServer struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port string `env:"PORT" env-default:"8080"`
}
type Database struct {
	Db  string `env:"DB" env-default:"mongo"`
	Url string `env:"DB_URL"`
}

// Load config from enviroment
// Throw an error if broker connection string is not setted
func LoadConfig() (Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.Database.Url == "" {
		log.Println("db string not setted")
	}

	return cfg, nil
}
