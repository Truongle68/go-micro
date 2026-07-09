package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

const (
	_defaultAddr           = ":80"
	_defaulReadTimeout     = 5 * time.Second
	_defaulWriteTimeout    = 5 * time.Second
	_defaulShutdownTimeout = 3 * time.Second
)

type Server struct {
	ctx context.Context
	eg  *errgroup.Group

	Engine     *gin.Engine
	httpServer *http.Server
	notify     chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	logger logger.Interface
}

func New(l logger.Interface, opts ...Option) *Server {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(1)

	s := &Server{
		ctx:             ctx,
		eg:              eg,
		Engine:          nil,
		notify:          make(chan error),
		address:         _defaultAddr,
		readTimeout:     _defaulReadTimeout,
		writeTimeout:    _defaulWriteTimeout,
		shutdownTimeout: _defaulShutdownTimeout,
		logger:          l,
	}

	for _, opt := range opts {
		opt(s)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	s.Engine = engine

	s.httpServer = &http.Server{
		Addr:         s.address,
		Handler:      s.Engine,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	return s
}

func (s *Server) Start() {
	s.eg.Go(func() error {
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
			close(s.notify)
			return err
		}
		return nil
	})

	s.logger.Info("http server - Server - Started")
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	var errs []error
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error("http server - Server - Shutdown - s.httpServer.Shutdown")
		errs = append(errs, err)
	}

	if err := s.eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		s.logger.Error("http server - Server - Shutdown - s.eg.Wait")
		errs = append(errs, err)
	}

	s.logger.Info("http server - Server - Shutdown")

	return errors.Join(errs...)
}
