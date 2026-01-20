package main

import (
	"github.com/scmbr/subscription-aggregator/internal/app"
	_ "github.com/scmbr/subscription-aggregator/internal/docs"
)

const configsDir = "configs"

// @title Subscription Aggregator API
// @version 1.0
// @description API for managing subscriptions
// @BasePath /api/v1

func main() {
	app.Run(configsDir)
}
