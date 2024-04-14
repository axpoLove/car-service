//go:generate swag init -pd
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaswdr/faker/v2"
	_ "github.com/lib/pq"

	_ "car-service/cmd/car-info-api/docs"
	"car-service/internal/config"
	server "car-service/internal/server/car-info"
	service "car-service/internal/service/car-info"
)

// @title car-info-api
// @version 1.0
// @description 'This is a car info api server'
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

	svc := service.NewService(faker.New())

	srv := server.NewServer(cfg.Server.CarInfoPort, cfg.Server.ReadHeaderTimeout, svc, logger)
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
