package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpServer
	Database
	Logger
}

type HttpServer struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port string `env:"PORT" env-default:"8080"`
}
type Database struct {
	Db       string `env:"DB" env-default:"mongo"`
	Name     string `env:"DB_NAME" env-default:"arcanery"`
	Url      string `env:"DB_URL"`
	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
}

type Logger struct {
	Log string `env:"LOG" env-default:"sentry"`
	Url string `env:"LOG_URL"`
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
