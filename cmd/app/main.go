package main

import (
	"github.com/scmbr/subscription-aggregator/internal/app"
)

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
