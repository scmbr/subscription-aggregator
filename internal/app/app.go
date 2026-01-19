package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/scmbr/subscription-aggregator/internal/config"
	"github.com/scmbr/subscription-aggregator/internal/delivery/http/handler"
	"github.com/scmbr/subscription-aggregator/internal/repository"
	"github.com/scmbr/subscription-aggregator/internal/server"
	"github.com/scmbr/subscription-aggregator/internal/service"
	"github.com/scmbr/subscription-aggregator/pkg/database/postgres"
	"github.com/scmbr/subscription-aggregator/pkg/logger"
)

func Run(configsDir string) {
	cfg, err := config.Init(configsDir)
	if err != nil {
		logger.Error("failed to initialize configs", err, map[string]interface{}{
			"configs_directory": configsDir,
		})
	}
	logger.Info("configs initialized successfully", map[string]interface{}{
		"http_port":     cfg.HTTP.Port,
		"postgres_port": cfg.Postgres.Port,
	})

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.Name,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		logger.Error("failed to connect to database", err, nil)
	}
	logger.Info("connected to database successfully", map[string]interface{}{
		"database is connected": db.DB.Ping() == nil,
	})
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)
	server := server.NewServer(cfg, handler.Init())
	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error occurred while running server", err, nil)
		}
	}()
	logger.Info("server started", nil)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	if err := server.Stop(ctx); err != nil {
		logger.Error("failed to stop server", err, nil)
	}
}
