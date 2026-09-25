package main

import (
	"context"
	"log"
	"order-service/config"
	"order-service/internal/domain"
	pgrepo "order-service/internal/repo/postgres"
	"order-service/internal/worker"
	pgtrans "order-service/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/postgres"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("starting outbox publisher service...")

	// init postgres connection & transactor
	pg, err := postgres.New(cfg.PG.Url)
	if err != nil {
		l.Fatal("initializing postgres: %v", err)
	}

	transactor := pgtrans.NewPostgresTransactor(pg.DB)

	// init outbox repo
	outboxRepo := pgrepo.NewOutboxRepo(pg.DB)

	pubConn := rabbitmq.NewConnection(cfg.RMQ.URL)
	if err := pubConn.Connect(); err != nil {
		l.Fatal("failed to connect rabbitmq: %v", err)
	}
	defer pubConn.Close()

	routes := map[string]rabbitmq.ExchangeBinding{
		domain.EventOrderPlaced: {Exchange: rabbitmq.ExchangeOrderPlaced, ExchangeType: rabbitmq.ExchangeTypeFanout},
	}

	// init RabbitMQ event publisher
	rmqPublisher, err := rabbitmq.NewPublisher(pubConn.Conn, true, routes)
	if err != nil {
		l.Warn("failed to initialize RabbitMQ publisher (events disabled): %v", err)
	}
	defer rmqPublisher.Close()

	// init worker with functional options
	outboxWorker := worker.NewOutBoxPublisherWorker(
		outboxRepo,
		transactor,
		rmqPublisher,
		l,
		worker.WithPollInterval(cfg.Outbox.PollInterval),
		worker.WithBatchSize(cfg.Outbox.BatchSize),
		worker.WithMaxRetries(cfg.Outbox.MaxRetries),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start worker in background
	go outboxWorker.Start(ctx)

	// wait for termination signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	sig := <-interrupt
	l.Info("outbox publisher received signal: %s, shutting down...", sig.String())

	cancel()
	outboxWorker.Stop()
	l.Info("outbox publisher stopped gracefully")
}
