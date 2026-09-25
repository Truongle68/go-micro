package config

import (
	"fmt"
	"time"

	env "github.com/TruongLe68/go-micro/pkg/config"
)

type (
	Config struct {
		HTTP     http
		PG       pg
		Services services
		JWT      jwt
		Redis    redis
		RMQ      rmq
		Outbox   outbox
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
		CatalogServiceGRPCAddr   string `env:"CATALOG_SERVICE_GRPC_ADDR" envDefault:"localhost:50051"`
		UserServiceGRPCAddr      string `env:"USER_SERVICE_GRPC_ADDR" envDefault:"localhost:50050"`
		InventoryServiceGRPCAddr string `env:"INVENTORY_SERVICE_GRPC_ADDR" envDefault:"localhost:50052"`
		CartServiceGRPCAddr      string `env:"CART_SERVICE_GRPC_ADDR" envDefault:"localhost:50053"`
	}

	jwt struct {
		PublicKey string `env:"PUBLIC_KEY,required"`
	}

	redis struct {
		Addr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
		DB       int    `env:"REDIS_DB" envDefault:"0"`
	}

	rmq struct {
		URL          string `env:"RMQ_URL,required"`
		Exchange     string `env:"RMQ_EXCHANGE" envDefault:"inventory.events"`
		ExchangeType string `env:"RMQ_EXCHANGE_TYPE" envDefault:"topic"`
	}

	outbox struct {
		PollInterval time.Duration `env:"OUTBOX_POLL_INTERVAL" envDefault:"5s"`
		BatchSize    int           `env:"OUTBOX_BATCH_SIZE" envDefault:"50"`
		MaxRetries   int           `env:"OUTBOX_MAX_RETRIES" envDefault:"5"`
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
