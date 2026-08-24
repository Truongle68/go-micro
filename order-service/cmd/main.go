package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"order-service/config"
	grpcclient "order-service/internal/client/grpc"
	"order-service/internal/delivery/http"
	v1 "order-service/internal/delivery/http/v1"
	pgrepo "order-service/internal/repo/postgres"
	"order-service/internal/usecase"
	pgtransactor "order-service/pkg/postgres"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/postgres"
	"github.com/TruongLe68/go-micro/pkg/redis"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	l := logger.New(cfg.Log.Level)

	// init postgres db
	pg, err := postgres.New(cfg.PG.Url)
	if err != nil {
		l.Fatal("failed to initialize postgres: %v", err)
	}
	defer pg.Close()
	transactor := pgtransactor.NewPostgresTransactor(pg.DB)

	// optional redis connection for token blacklist checking
	redisClient, err := redis.New("localhost:6379", "secretredispass", 0)
	if err != nil {
		l.Warn("redis connection optional warning: %v", err)
	}
	var cache *redismanager.IdentityCache
	if redisClient != nil {
		defer redisClient.Close()
		cache = redismanager.NewIdentityCache(redisClient.Client)
	}

	// init JWT verifier
	jwtVerifier, err := jwtmanager.NewVerifier(cfg.JWT.PublicKey)
	if err != nil {
		l.Fatal("failed to initialize jwt verifier: %v", err)
	}

	// init gRPC clients
	cartClient, err := grpcclient.NewCartGRPCClient(cfg.Services.CartServiceGRPCAddr)
	if err != nil {
		l.Fatal("failed to initialize cart gRPC client: %v", err)
	}
	catalogClient, err := grpcclient.NewCatalogGRPCClient(cfg.Services.CatalogServiceGRPCAddr)
	if err != nil {
		l.Fatal("failed to initialize catalog gRPC client: %v", err)
	}

	// init repo & usecase
	orderRepo := pgrepo.NewOrderRepo(pg.DB)
	orderUC := usecase.NewOrderUC(orderRepo, cartClient, catalogClient, transactor, l)

	// init http server & router
	server := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))

	var blacklistCacher redismanager.BlacklistCacher
	if cache != nil {
		blacklistCacher = cache.TokenBlacklist
	}

	v1Deps := v1.NewDependencies(orderUC, jwtVerifier, blacklistCacher, l)
	http.NewRouter(server.Engine, v1Deps)

	// start http server
	server.Start()

	// graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Info("app - Run - interrupt: %s", sig.String())
	case err = <-server.Notify():
		l.Error(fmt.Errorf("app - Run - server.Notify: %w", err))
	}

	if err := server.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - server.Shutdown: %w", err))
	}
}
