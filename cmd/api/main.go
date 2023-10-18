package main

import (
	"context"
	"log"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app/config"
)

func main() {
	ctx := context.Background()
	cfg := config.New()
	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to init app: %v", err)
	}
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
