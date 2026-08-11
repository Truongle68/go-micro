package app

import (
	"catalog-service/config"
	repo "catalog-service/internal/repo/mongodb"
	"catalog-service/internal/usecase"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	httpr "catalog-service/internal/delivery/http"
	v1 "catalog-service/internal/delivery/http/v1"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/mongo"
	"github.com/TruongLe68/go-micro/pkg/redis"
	mongodrv "go.mongodb.org/mongo-driver/v2/mongo"
)

type servers struct {
	http *httpserver.Server
}

type useCases struct {
	product  usecase.ProductUsecase
	category usecase.CategoryUsecase
}

func initUsecases(db *mongodrv.Database) useCases {
	productRepo := repo.NewProductRepo(db)
	if err := productRepo.EnsureIndexes(context.Background()); err != nil {
		log.Printf("failed to ensure product indexes: %v", err)
	}
	categoryRepo := repo.NewCategoryRepo(db)
	productUC := usecase.NewProductUC(productRepo, categoryRepo)
	categoryUC := usecase.NewCategoryUC(categoryRepo)

	return useCases{
		product:  productUC,
		category: categoryUC,
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

	return &servers{
		http: http,
	}
}

func (s *servers) startServer() {
	s.http.Start()
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
	}
}

func (s *servers) shutdownServers(l logger.Interface) {
	if err := s.http.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - s.http.Shutdown: %w", err))
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
