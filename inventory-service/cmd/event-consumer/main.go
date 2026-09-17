package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"inventory-service/config"
	grpcclient "inventory-service/internal/client/grpc"
	"inventory-service/internal/domain"
	pgrepo "inventory-service/internal/repo/postgres"
	"inventory-service/internal/usecase"
	"inventory-service/internal/worker"
	invpg "inventory-service/pkg/postgres"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/postgres"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("starting inventory event consumer...")

	// Infrastructure: DB
	pg, err := postgres.New(cfg.PG.Url)
	if err != nil {
		l.Fatal("postgres: %v", err)
	}
	defer pg.Close()
	transactor := invpg.NewPostgresTransactor(pg.DB)

	// init RabbitMQ event publisher
	var eventPublisher usecase.EventPublisher
	rmqPublisher, err := rabbitmq.NewPublisher(cfg.RMQ.URL, cfg.RMQ.Exchange)
	if err != nil {
		l.Warn("failed to initialize RabbitMQ publisher (events disabled): %v", err)
	} else {
		eventPublisher = rmqPublisher
		defer rmqPublisher.Close()
	}

	// init repos
	warehouseRepo := pgrepo.NewWarehouseRepo(pg.DB)
	stockLevelRepo := pgrepo.NewStockLevelRepo(pg.DB)
	stockReservationRepo := pgrepo.NewStockReservationRepo(pg.DB)
	stockMovementRepo := pgrepo.NewStockMovementRepo(pg.DB)

	// init clients
	catalogClient, err := grpcclient.NewCatalogGRPCClient(cfg.Services.CatalogServiceAddr)
	if err != nil {
		l.Fatal("failed to initialize catalog gRPC client: %v", err)
	}

	// init usecases
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

	// register handler
	handler := worker.NewGoodsReceivedHandler(stockUC, l)

	exchange := cfg.RMQ.Exchange
	if exchange == "" {
		exchange = "inventory.events"
	}

	queueName := cfg.RMQ.QueueName
	if queueName == "" {
		queueName = "inventory.goods-received.q"
	}

	routingKey := cfg.RMQ.RoutingKey
	if routingKey == "" {
		routingKey = domain.EventGoodsReceived
	}
	// Create consumer
	c, err := rabbitmq.NewConsumer(
		cfg.RMQ.URL,
		exchange,
		queueName,
		routingKey,
		handler.Handle,
	)
	if err != nil {
		l.Fatal("rabbitmq consumer: %v", err)
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := c.Start(ctx); err != nil && ctx.Err() == nil {
			l.Error("consumer stopped: %v", err)
			cancel()
		}
	}()

	l.Info("inventory event consumer started")

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	sig := <-interrupt
	l.Info("consumer received signal: %s, shutting down...", sig.String())
	cancel()
	l.Info("inventory event consumer stopped")
}
