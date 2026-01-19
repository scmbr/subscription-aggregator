package app

import (
	"github.com/scmbr/subscription-aggregator/internal/config"
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
}
