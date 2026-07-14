package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"user-service/config"
	"user-service/internal/delivery/http"
	v1 "user-service/internal/delivery/http/v1"
	repo "user-service/internal/repo/postgres"
	"user-service/internal/usecase"
	"user-service/pkg/jwt"
	userredis "user-service/pkg/redis"

	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/postgres"
	"github.com/TruongLe68/go-micro/pkg/redis"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	l := logger.New(cfg.Log.Level)

	// init db conn
	pg, err := postgres.New(cfg.PG.Url)
	if err != nil {
		l.Fatal("failed to initialize postgres: %v", err)
	}
	defer pg.Close()

	// init redis conn
	redisClient, err := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("failed to initialize redis: %v", err)
	}
	defer redisClient.Close()

	// init jwt
	jwtManager := jwt.New(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	// init repo + usecase
	userRepo := repo.NewUserRepo(pg.DB)

	blacklist := userredis.NewTokenBlacklist(redisClient.Client)
	profile := userredis.NewProfileCache(redisClient.Client)
	otpCache := userredis.NewOTPCache(redisClient.Client)

	authUC := usecase.NewAuthUC(userRepo, jwtManager, blacklist, profile, otpCache)
	userUC := usecase.NewUserUC(userRepo, profile)

	// init server, setup routers
	httpserver := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))

	v1Deps := v1.NewDependencies(authUC, userUC, jwtManager, redisClient.Client, l)
	http.NewRouter(httpserver.Engine, v1Deps)

	// start server
	httpserver.Start()

	// shutdown server
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Info("app - Run - interrupt: %s", sig.String())
	case err = <-httpserver.Notify():
		l.Error(fmt.Errorf("app - Run - s.http.Notify: %w", err))
	}

	if err := httpserver.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - s.http.Shutdown: %w", err))
	}
}
