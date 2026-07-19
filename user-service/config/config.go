package config

import (
	"fmt"
	"time"

	env "github.com/TruongLe68/go-micro/pkg/config"
)

type (
	Config struct {
		HTTP  http
		GRPC  grpc
		PG    pg
		JWT   jwt
		Redis redis
		Log   log
		Email email
	}

	http struct {
		Port    string `env:"HTTP_PORT" envDefault:"4000"`
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
		Port     string `env:"GRPC_PORT" envDefault:"50050"`
		Services services
	}

	services struct {
		RewardServiceAddr string `env:"REWARD_SERVICE_ADDR" envDefault:"50051"`
	}

	pg struct {
		Url string `env:"DB_URL,required"`
	}

	jwt struct {
		AccessSecret  string        `env:"JWT_ACC_SECRET,required"`
		RefreshSecret string        `env:"JWT_REF_SECRET,required"`
		AccessExpiry  time.Duration `env:"JWT_ACC_EXPIRY" envDefault:"15m"`
		RefreshExpiry time.Duration `env:"JWT_REF_EXPIRY" envDefault:"168h"`
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
