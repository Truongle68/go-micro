package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"user-service/config"

	"github.com/TruongLe68/go-micro/pkg/grpcserver"
	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
)

type servers struct {
	http *httpserver.Server
	grpc *grpcserver.Server
}

func initUsecases() {}

func initServers(cfg *config.Config, l logger.Interface) *servers {
	http := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))

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

	s := initServers(cfg, l)
	// start server
	s.startServer()
	// wait for server shutdown
	s.waitForShutdown(l)
}
