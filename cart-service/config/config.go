package config

import (
	"fmt"
	"time"

	env "github.com/TruongLe68/go-micro/pkg/config"
)

type (
	Config struct {
		HTTP  HTTP
		GRPC  GRPC
		Services Services
		Redis Redis
		JWT   JWT
		Log   Log
		Cart  Cart
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" envDefault:"4003"`
	}

	GRPC struct {
		Port string `env:"GRPC_PORT" envDefault:"50053"`
	}

	Services struct {
		CatalogServiceGRPCAddr string `env:"CATALOG_SERVICE_GRPC_ADDR" envDefault:"localhost:50050"`
	}

	Redis struct {
		Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}

	JWT struct {
		PublicKey string `env:"PUBLIC_KEY,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug"`
	}

	Cart struct {
		TTL time.Duration `env:"CART_TTL" envDefault:"720h"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Load(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
