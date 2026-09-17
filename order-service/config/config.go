package config

import (
	"fmt"

	env "github.com/TruongLe68/go-micro/pkg/config"
)

type (
	Config struct {
		HTTP     http
		PG       pg
		Services services
		JWT      jwt
		Redis    redis
		Log      log
	}

	http struct {
		Port string `env:"HTTP_PORT" envDefault:"4004"`
	}

	pg struct {
		Url string `env:"DB_URL,required"`
	}

	services struct {
		CartServiceURL           string `env:"CART_SERVICE_URL" envDefault:"http://localhost:4003"`
		CartServiceGRPCAddr      string `env:"CART_SERVICE_GRPC_ADDR" envDefault:"localhost:50053"`
		CatalogServiceGRPCAddr   string `env:"CATALOG_SERVICE_GRPC_ADDR" envDefault:"localhost:50050"`
		InventoryServiceGRPCAddr string `env:"INVENTORY_SERVICE_GRPC_ADDR" envDefault:"localhost:50052"`
	}

	jwt struct {
		PublicKey string `env:"PUBLIC_KEY,required"`
	}

	redis struct {
		Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}

	log struct {
		Level string `env:"LOG_LEVEL" envDefault:"debug"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Load(cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
