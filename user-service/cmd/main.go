package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"user-service/config"
	"user-service/internal/delivery/http"
	v1 "user-service/internal/delivery/http/v1"
	repo "user-service/internal/repo/postgres"
	"user-service/internal/usecase"
	"user-service/pkg/jwt"
	"user-service/pkg/mailer"
	userpg "user-service/pkg/postgres"
	userredis "user-service/pkg/redis"
	"user-service/pkg/sms"
	"user-service/worker"

	"github.com/TruongLe68/go-micro/pkg/httpserver"
	"github.com/TruongLe68/go-micro/pkg/logger"
	"github.com/TruongLe68/go-micro/pkg/postgres"
	"github.com/TruongLe68/go-micro/pkg/redis"
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
	transactor := userpg.NewPostgresTransactor(pg.DB)

	// init redis conn
	redisClient, err := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		l.Fatal("failed to initialize redis: %v", err)
	}
	defer redisClient.Close()

	// init jwt
	jwtManager, err := jwt.New(cfg.JWT.PrivateKey, cfg.JWT.PublicKey, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry, cfg.JWT.ActionExpiry)
	if err != nil {
		l.Fatal("failed to initialize jwt: %v", err)
	}

	// init repo, cache
	userRepo := repo.NewUserRepo(pg.DB)

	cache := userredis.NewIdentityCache(redisClient.Client)

	// init mailer, worker
	emailMailer := mailer.New(cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.SMTPUser, cfg.Email.SMTPPassword, cfg.Email.SenderEmail, l)
	emailWorker := worker.NewEmailWorker(emailMailer, l)
	smsClient := sms.NewMockSMS(l)

	smsWorker := worker.NewSMSWorker(smsClient, l)
	otpDispatcher := worker.NewCompositeOTPDispatcher(emailWorker, smsWorker)

	// init usecase
	authUC := usecase.NewAuthUC(userRepo, jwtManager, cache, transactor, emailWorker, otpDispatcher, l, cfg.HTTP.BaseURL)
	userUC := usecase.NewUserUC(userRepo, jwtManager, cache, transactor, emailWorker, l, cfg.HTTP.BaseURL)

	// init server, setup routers
	httpserver := httpserver.New(l, httpserver.Port(cfg.HTTP.Port))

	v1Deps := v1.NewDependencies(authUC, userUC, jwtManager, cache.TokenBlacklist, l)
	http.NewRouter(httpserver.Engine, v1Deps)

	// start server
	httpserver.Start()

	// shutdown server
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		l.Info("app - Run - interrupt: %s", sig.String())
	case err = <-httpserver.Notify():
		l.Error(fmt.Errorf("app - Run - s.http.Notify: %w", err))
	}

	if err := httpserver.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - s.http.Shutdown: %w", err))
	}
}
