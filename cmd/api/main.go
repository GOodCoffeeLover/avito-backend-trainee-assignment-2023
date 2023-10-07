package main

import (
	"log"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app/config"
)

func main() {
	cfg := &config.Config{
		Port: 7001,
	}
	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
