package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"cart-service/config"
	"cart-service/internal/client/grpc"
	grpcv1 "cart-service/internal/delivery/grpc/v1"
	"cart-service/internal/delivery/http"
	v1 "cart-service/internal/delivery/http/v1"
	redisrepo "cart-service/internal/repo/redis"
	"cart-service/internal/usecase"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	cartv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/cart/v1"
	"github.com/TruongLe68/go-micro/pkg/grpcserver"
	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/redis"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	l := logger.New(cfg.Log.Level)

	// init redis connection
	redisClient, err := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("failed to initialize redis: %v", err)
	}
	defer redisClient.Close()

	// init JWT verifier
	jwtVerifier, err := jwtmanager.NewVerifier(cfg.JWT.PublicKey)
	if err != nil {
		l.Fatal("failed to initialize jwt verifier: %v", err)
	}

	cache := redismanager.NewIdentityCache(redisClient.Client)

	// init gRPC client
	catalogClient, err := grpc.NewCatalogGRPCClient(cfg.Services.CatalogServiceGRPCAddr)

	// init repo & usecase
	cartRepo := redisrepo.NewCartRepo(redisClient.Client, cfg.Cart.TTL)
	cartUC := usecase.NewCartUC(cartRepo, catalogClient, l)

	// init http server & router
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))

	v1Deps := v1.NewDependencies(cartUC, jwtVerifier, cache.TokenBlacklist, l)
	http.NewRouter(httpServer.Engine, v1Deps)

	// init grpc server
	grpcServer := grpcserver.New(l, grpcserver.Port(cfg.GRPC.Port))
	cartGrpcServer := grpcv1.NewCartServer(cartUC, l)
	cartv1.RegisterCartServiceServer(grpcServer.App, cartGrpcServer)

	// start servers
	httpServer.Start()
	grpcServer.Start()

	// graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Info("app - Run - interrupt: %s", sig.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	if err := httpServer.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
	if err := grpcServer.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}
}
