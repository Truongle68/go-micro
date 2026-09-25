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
	"golang.org/x/sync/errgroup"
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

	// init rabbitmq conn
	conn := rabbitmq.NewConnection(cfg.RMQ.URL)
	if err := conn.Connect(); err != nil {
		l.Fatal("failed to connect rabbitmq: %v", err)
	}
	defer conn.Close()

	subs := worker.BuildSubscriptions(stockUC, l)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, gCtx := errgroup.WithContext(ctx)
	consumers := make([]*rabbitmq.Consumer, 0, len(subs))

	for _, s := range subs {
		s := s
		consumer, err := s.Register(conn.Conn)
		if err != nil {
			l.Fatal("failed to init consumer: %v", err)
		}

		consumers = append(consumers, consumer)
		g.Go(func() error {
			l.Info("consumer %s started (queue=%s)", s.Name, s.Queue)
			return consumer.Start(gCtx)
		})
	}

	l.Info("inventory event consumer started")

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Info("consumer received signal: %s, shutting down...", sig.String())
	case <-gCtx.Done():
		l.Error("a consumer stopped unexpectedly, shutting down...")
	}

	cancel()

	for _, c := range consumers {
		if err := c.Close(); err != nil {
			l.Error("error closing consumer: %v", err)
		}
	}

	if err := g.Wait(); err != nil && ctx.Err() == nil {
		l.Error("consumer group exited with error: %v", err)
	}

	l.Info("inventory event consumer stopped")
}
