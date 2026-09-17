package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"inventory-service/config"
	pgrepo "inventory-service/internal/repo/postgres"
	"inventory-service/internal/worker"
	invpg "inventory-service/pkg/postgres"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/postgres"
	"github.com/TruongLe68/go-micro/pkg/rabbitmq"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l := logger.New(cfg.Log.Level)
	l.Info("starting outbox publisher service...")

	// init postgres connection & transactor
	pg, err := postgres.New(cfg.PG.Url)
	if err != nil {
		l.Fatal("failed to initialize postgres: %v", err)
	}
	defer pg.Close()
	transactor := invpg.NewPostgresTransactor(pg.DB)

	// init outbox repo
	outboxRepo := pgrepo.NewOutboxRepo(pg.DB)

	// init RabbitMQ event publisher
	exchangeName := cfg.RMQ.Exchange
	if exchangeName == "" {
		exchangeName = "inventory.events"
	}
	exchangeType := cfg.RMQ.ExchangeType
	if exchangeType == "" {
		exchangeType = "topic"
	}
	rmqPublisher, err := rabbitmq.NewPublisher(cfg.RMQ.URL, exchangeName, rabbitmq.ExchangeType(exchangeType))
	if err != nil {
		l.Fatal("failed to initialize RabbitMQ publisher: %v", err)
	}
	defer rmqPublisher.Close()

	// init worker with functional options
	outboxWorker := worker.NewOutboxPublisherWorker(
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
