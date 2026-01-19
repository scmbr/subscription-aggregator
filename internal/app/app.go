package app

import (
	"github.com/scmbr/subscription-aggregator/internal/config"
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
}
