package config

import (
	"fmt"

	env "github.com/TruongLe68/go-micro/pkg/config"
)

type (
	Config struct {
		HTTP     HTTP
		PG       PG
		Services Services
		JWT      JWT
		Log      Log
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" envDefault:"4004"`
	}

	PG struct {
		Url string `env:"DB_URL,required"`
	}

	Services struct {
		CartServiceURL             string `env:"CART_SERVICE_URL" envDefault:"http://localhost:4003"`
		CartServiceGRPCAddr        string `env:"CART_SERVICE_GRPC_ADDR" envDefault:"localhost:50053"`
		CatalogServiceGRPCAddr     string `env:"CATALOG_SERVICE_GRPC_ADDR" envDefault:"localhost:50050"`
		InventoryServiceGRPCAddr   string `env:"INVENTORY_SERVICE_GRPC_ADDR" envDefault:"localhost:50052"`
	}

	JWT struct {
		PublicKey string `env:"PUBLIC_KEY,required"`
	}

	Log struct {
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
