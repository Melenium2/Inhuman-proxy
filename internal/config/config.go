package config

import (
	"github.com/caarlos0/env/v6"
	"net/url"
)

type Config struct {
	Port          int       `env:"PORT" envDefault:"19000"`
	Servers       []url.URL `env:"SERVERS"`
	StorageConfig struct {
		Addr     string `env:"REDIS_ADDR" envDefault:"192.168.99.100:6379"`
		Username string `env:"REDIS_USERNAME,required"`
		Password string `env:"REDIS_PASSWORD,required"`
		Database int    `env:"REDIS_DATABASE" envDefault:"0"`
	}
}

type StorageConfig struct {
	Addr     string `env:"REDIS_ADDR" envDefault:"192.168.99.100:6379"`
	Username string `env:"REDIS_USERNAME,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	Database int    `env:"REDIS_DATABASE" envDefault:"0"`
}

func New() (Config, error) {
	cfg := Config{}
	if err := env.Parse(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
