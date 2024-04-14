//go:generate swag init -pd
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocraft/dbr/v2"
	_ "github.com/lib/pq"

	_ "car-service/cmd/car-api/docs"
	carinfoclient "car-service/internal/client/car-info"
	"car-service/internal/config"
	repository "car-service/internal/repository/car"
	server "car-service/internal/server/car"
	service "car-service/internal/service/car"
)

// @title car-api
// @version 1.0
// @description 'This is a car api server'
// @BasePath /
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		return
	}
	db, err := dbr.Open("postgres", cfg.Postgres.ConnString(), nil)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error("failed to close db connection", "error", err)
		}
	}()
	err = db.Ping()
	if err != nil {
		logger.Error("failed to ping postgres", "error", err)
		return
	}

	svc := service.NewService(
		repository.NewRepository(db),
		carinfoclient.NewClient(cfg.CarInfoService.URL, &http.Client{
			Timeout: cfg.CarInfoService.RequestTimeout,
		}),
	)

	srv := server.NewServer(cfg.Server.CarPort, cfg.Server.ReadHeaderTimeout, svc, logger)
	serverErrors := make(chan error)
	go func() {
		err := srv.Start()
		if err != nil {
			serverErrors <- err
		}
	}()
	logger.Info("server started")

	select {
	case <-ctx.Done():
	case err := <-serverErrors:
		logger.Error("server failed", "error", err)
		return
	}
	logger.Info("server stopped")
}
