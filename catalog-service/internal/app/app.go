package app

import (
	"catalog-service/config"
	repo "catalog-service/internal/repo/mongodb"
	"catalog-service/internal/usecase"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	httpr "catalog-service/internal/delivery/http"
	v1 "catalog-service/internal/delivery/http/v1"

	"github.com/TruongLe68/go-micro/pkg/grpcserver"
	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/mongo"
	mongodrv "go.mongodb.org/mongo-driver/v2/mongo"
)

type servers struct {
	http *httpserver.Server
	grpc *grpcserver.Server
}

type useCases struct {
	product  usecase.ProductUsecase
	category usecase.CategoryUsecase
}

func initUsecases(db *mongodrv.Database) useCases {
	productRepo := repo.NewProductRepo(db)
	categoryRepo := repo.NewCategoryRepo(db)
	productUC := usecase.NewProductUC(productRepo)
	categoryUC := usecase.NewCategoryUC(categoryRepo)

	return useCases{
		product:  productUC,
		category: categoryUC,
	}
}

func initServers(l logger.Interface, uc useCases, cfg *config.Config) *servers {
	http := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))
	deps := v1.Dependencies{
		Product:   uc.product,
		Category:  uc.category,
		Logger:    l,
		JWTSecret: cfg.JWT.AccessSecret,
	}
	httpr.NewRouter(http.Engine, deps)

	grpc := grpcserver.New(l, grpcserver.Port(cfg.GRPC.Port))

	return &servers{
		http: http,
		grpc: grpc,
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
	case err = <-s.grpc.Notify():
		l.Error(fmt.Errorf("app - Run - s.http.Notify: %w", err))
	}

	s.shutdownServers(l)
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
	// init usecase
	uc := initUsecases(m.Database)
	// init server
	s := initServers(l, uc, cfg)
	// start server
	s.startServer()
	// wait for server shutdown
	s.waitForShutdown(l)
}
