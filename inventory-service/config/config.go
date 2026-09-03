package config

import (
	"fmt"

	env "github.com/TruongLe68/go-micro/pkg/config"
)

type (
	Config struct {
		HTTP     http
		GRPC     grpc
		Services services
		RMQ      rmq
		PG       pg
		JWT      jwt
		Redis    redis
		Log      log
		Email    email
	}

	http struct {
		Port    string `env:"HTTP_PORT" envDefault:"4002"`
		BaseURL string `env:"BASE_URL" envDefault:"http://localhost:3000"`
	}

	email struct {
		SMTPHost     string `env:"SMTP_HOST" envDefault:""`
		SMTPPort     int    `env:"SMTP_PORT" envDefault:"587"`
		SMTPUser     string `env:"SMTP_USER" envDefault:""`
		SMTPPassword string `env:"SMTP_PASSWORD" envDefault:""`
		SenderEmail  string `env:"SENDER_EMAIL" envDefault:"no-reply@example.com"`
	}

	grpc struct {
		Port string `env:"GRPC_PORT" envDefault:"50052"`
	}

	rmq struct {
		ServerExchange string `env:"RMQ_RPC_SERVER,required"`
		ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
		URL            string `env:"RMQ_URL,required"`
	}

	services struct {
		CatalogServiceAddr string `env:"CATALOG_SERVICE_ADDR" envDefault:"localhost:50050"`
	}

	pg struct {
		Url string `env:"DB_URL,required"`
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
