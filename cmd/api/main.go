package main

import (
	"context"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app/config"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed get config")
	}

	app, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to init app")
	}

	if err := app.Run(ctx); err != nil {
		log.Fatal().Err(err).Msgf("Failed to run app")
	}
}
