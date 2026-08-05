package config

import (
	"fmt"

	env "github.com/TruongLe68/go-micro/pkg/config"
)

type (
	Config struct {
		HTTP    http
		GRPC    grpc
		MongoDB mongodb
		Log     log
		JWT     jwt
		Redis   redis
	}

	http struct {
		Port string `env:"HTTP_PORT" envDefault:"4001"`
	}

	grpc struct {
		Port     string `env:"GRPC_PORT" envDefault:"50050"`
		Services services
	}

	services struct {
		UserServiceAddr string `env:"USER_SERVICE_ADDR" envDefault:"50050"`
	}

	mongodb struct {
		Url    string `env:"DB_URL,required"`
		DBName string `env:"DB_NAME,required"`
	}

	log struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug"`
	}

	jwt struct {
		PublicKey string `env:"PUBLIC_KEY,required"`
	}

	redis struct {
		Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Load(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
