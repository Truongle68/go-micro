package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"golang.org/x/sync/errgroup"
	pbgrpc "google.golang.org/grpc"
)

const _defaultAddr = ":80"

type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	App        *pbgrpc.Server
	notify     chan error
	address    string
	serverOpts []pbgrpc.ServerOption

	logger logger.Interface
}

func New(l logger.Interface, opts ...Option) *Server {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(1)

	s := &Server{
		ctx:     ctx,
		eg:      eg,
		notify:  make(chan error, 1),
		address: _defaultAddr,
		logger:  l,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.App = pbgrpc.NewServer(s.serverOpts...)

	return s
}

func (s *Server) Start() {
	s.eg.Go(func() error {
		var lc net.ListenConfig
		ln, err := lc.Listen(s.ctx, "tcp", s.address)
		if err != nil {
			s.notify <- err
			close(s.notify)
			return err
		}

		err = s.App.Serve(ln)
		if err != nil {
			s.notify <- err
			close(s.notify)
			return err
		}
		return nil
	})

	s.logger.Info(fmt.Sprintf("grpc server - Server - Started at port: %s", s.address))
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	s.App.GracefulStop()

	if err := s.eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error("grpc server - Shutdown - s.eg.Wait")
		return err
	}

	s.logger.Info("grpc server - Server - Shutdown")
	return nil
}
