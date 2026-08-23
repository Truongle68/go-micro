package app

import (
	"catalog-service/config"
	"catalog-service/internal/repo/mongorepo"
	"catalog-service/internal/usecase"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	grpcv1 "catalog-service/internal/delivery/grpc/v1"
	httpr "catalog-service/internal/delivery/http"
	v1 "catalog-service/internal/delivery/http/v1"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	catalogv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/catalog/v1"
	"github.com/TruongLe68/go-micro/pkg/grpcserver"
	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/mongo"
	"github.com/TruongLe68/go-micro/pkg/redis"
	mongodrv "go.mongodb.org/mongo-driver/v2/mongo"
)

type servers struct {
	http *httpserver.Server
	grpc *grpcserver.Server
}

type useCases struct {
	product     v1.ProductUsecase
	grpcProduct grpcv1.ProductGRPCUsecase
	category    v1.CategoryUsecase
}

func initUsecases(db *mongodrv.Database) useCases {
	productRepo := mongorepo.NewProductRepo(db)
	if err := productRepo.EnsureIndexes(context.Background()); err != nil {
		log.Printf("failed to ensure product indexes: %v", err)
	}
	categoryRepo := mongorepo.NewCategoryRepo(db)
	productUC := usecase.NewProductUC(productRepo, categoryRepo)
	categoryUC := usecase.NewCategoryUC(categoryRepo)

	return useCases{
		product:     productUC,
		grpcProduct: productUC,
		category:    categoryUC,
	}
}

func initServers(l logger.Interface, uc useCases, cfg *config.Config, v jwtmanager.JWTManager, cache redismanager.IdentityCacher) *servers {
	http := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	deps := v1.Dependencies{
		Product:  uc.product,
		Category: uc.category,
		Logger:   l,
		Verifier: v,
		Cache:    cache,
	}
	httpr.NewRouter(http.Engine, deps)

	grpcServer := grpcserver.New(l, grpcserver.Port(cfg.GRPC.Port))
	productServer := grpcv1.NewProductServer(uc.grpcProduct, l)
	catalogv1.RegisterProductServiceServer(grpcServer.App, productServer)

	return &servers{
		http: http,
		grpc: grpcServer,
	}
}

func (s *servers) startServer() {
	s.http.Start()
	s.grpc.Start()
}

func (s *servers) waitForShutdown(l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - interrupt: %s", sig.String())
	case err = <-s.http.Notify():
		l.Error(fmt.Errorf("app - Run - s.http.Notify: %w", err))
		s.shutdownServers(l)
	case err = <-s.grpc.Notify():
		l.Error(fmt.Errorf("app - Run - s.grpc.Notify: %w", err))
		s.shutdownServers(l)
	}
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - s.http.Shutdown: %w", err))
	}
	if err := s.grpc.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - s.grpc.Shutdown: %w", err))
	}
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	// init db, repo
	m, err := mongo.New(cfg.MongoDB.Url, cfg.MongoDB.DBName)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed to init mongodb: %w", err))
	}
	defer m.Close()

	verifier, err := jwtmanager.NewVerifier(cfg.JWT.PublicKey)
	if err != nil {
		l.Fatal("failed to initialize jwt: %v", err)
	}

	red, err := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("failed to initialize redis: %v", err)
	}
	defer red.Close()

	identityCache := redismanager.NewIdentityCache(red.Client)
	// init usecase
	uc := initUsecases(m.Database)
	// init server
	s := initServers(l, uc, cfg, verifier, identityCache)
	// start server
	s.startServer()
	// wait for server shutdown
	s.waitForShutdown(l)
}
