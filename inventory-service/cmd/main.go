package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"inventory-service/config"
	grpcclient "inventory-service/internal/client/grpc"
	grpcv1 "inventory-service/internal/delivery/grpc/v1"
	"inventory-service/internal/delivery/http"
	v1 "inventory-service/internal/delivery/http/v1"
	"inventory-service/internal/domain"
	pgrepo "inventory-service/internal/repo/postgres"
	"inventory-service/internal/usecase"
	invpg "inventory-service/pkg/postgres"

	"github.com/GoProOrg/core-go-pkg/jwtmanager"
	redismanager "github.com/GoProOrg/core-go-pkg/redismanager/identity"
	inventoryv1 "github.com/TruongLe68/go-micro/pkg/gen/proto/go/inventory/v1"
	"github.com/TruongLe68/go-micro/pkg/grpcmw"
	"github.com/TruongLe68/go-micro/pkg/grpcserver"
	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/postgres"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
	"github.com/TruongLe68/go-micro/pkg/redis"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	l := logger.New(cfg.Log.Level)

	// init db conn, transactor
	pg, err := postgres.New(cfg.PG.Url)
	if err != nil {
		l.Fatal("failed to initialize postgres: %v", err)
	}
	defer pg.Close()
	transactor := invpg.NewPostgresTransactor(pg.DB)

	// init redis conn
	redisClient, err := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("failed to initialize redis: %v", err)
	}
	defer redisClient.Close()

	// init jwt
	jwtVerifier, err := jwtmanager.NewVerifier(cfg.JWT.PublicKey)
	if err != nil {
		l.Fatal("failed to initialize jwt verifier: %v", err)
	}

	// init cache
	cache := redismanager.NewIdentityCache(redisClient.Client)

	pubConn := rabbitmq.NewConnection(cfg.RMQ.URL)
	if err := pubConn.Connect(); err != nil {
		l.Fatal("failed to connect rabbitmq: %v", err)
	}
	defer pubConn.Close()

	routes := map[string]rabbitmq.ExchangeBinding{
		domain.EventGoodsReceived: {Exchange: rabbitmq.ExchangeInventory, ExchangeType: rabbitmq.ExchangeTypeTopic},
		domain.EventOrderPlaced:   {Exchange: rabbitmq.ExchangeOrderPlaced, ExchangeType: rabbitmq.ExchangeTypeFanout},
	}

	// init RabbitMQ event publisher
	var eventPublisher usecase.EventPublisher
	rmqPublisher, err := rabbitmq.NewPublisher(pubConn.Conn, true, routes)
	if err != nil {
		l.Warn("failed to initialize RabbitMQ publisher (events disabled): %v", err)
	} else {
		eventPublisher = rmqPublisher
		defer rmqPublisher.Close()
	}


	// init repos
	supplierRepo := pgrepo.NewSupplierRepo(pg.DB)
	supplyRepo := pgrepo.NewSupplyRepo(pg.DB)
	warehouseRepo := pgrepo.NewWarehouseRepo(pg.DB)
	purchaseOrderRepo := pgrepo.NewPurchaseOrderRepo(pg.DB)
	stockLevelRepo := pgrepo.NewStockLevelRepo(pg.DB)
	stockReservationRepo := pgrepo.NewStockReservationRepo(pg.DB)
	stockMovementRepo := pgrepo.NewStockMovementRepo(pg.DB)
	outboxRepo := pgrepo.NewOutboxRepo(pg.DB)

	// init clients
	catalogClient, err := grpcclient.NewCatalogGRPCClient(cfg.Services.CatalogServiceAddr)
	if err != nil {
		l.Fatal("failed to initialize catalog gRPC client: %v", err)
	}

	// init usecases
	warehouseUC := usecase.NewWarehouseUC(warehouseRepo)
	supplierUC := usecase.NewSupplierUC(supplierRepo, l)
	supplyUC := usecase.NewSupplyUC(supplierRepo, supplyRepo, catalogClient)
	purchaseOrderUC := usecase.NewPurchaseOrderUC(
		purchaseOrderRepo,
		supplierRepo,
		warehouseRepo,
		stockLevelRepo,
		stockMovementRepo,
		catalogClient,
		transactor,
		outboxRepo,
		l,
	)
	stockUC := usecase.NewStockUC(
		stockLevelRepo,
		warehouseRepo,
		stockReservationRepo,
		stockMovementRepo,
		catalogClient,
		transactor,
		eventPublisher,
		l,
	)

	// init HTTP server, setup routers
	server := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))

	v1Deps := v1.NewDependencies(supplierUC, supplyUC, purchaseOrderUC, warehouseUC, stockUC, jwtVerifier, cache.TokenBlacklist, l)
	http.NewRouter(server.Engine, v1Deps)

	// init gRPC server
	grpcSrv := grpcserver.New(l, grpcserver.Port(cfg.GRPC.Port), grpcserver.ServerOptions(grpc.UnaryInterceptor(grpcmw.LoggingInterceptor)))
	inventoryServer := grpcv1.NewInventoryServer(stockUC, l)
	inventoryv1.RegisterInventoryServiceServer(grpcSrv.App, inventoryServer)

	// start servers
	server.Start()
	grpcSrv.Start()

	// shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Info("app - Run - interrupt: %s", sig.String())
	case err = <-server.Notify():
		l.Error(fmt.Errorf("app - Run - s.http.Notify: %w", err))
	case err = <-grpcSrv.Notify():
		l.Error(fmt.Errorf("app - Run - s.grpc.Notify: %w", err))
	}

	if err := server.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - s.http.Shutdown: %w", err))
	}
	if err := grpcSrv.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - s.grpc.Shutdown: %w", err))
	}
}
