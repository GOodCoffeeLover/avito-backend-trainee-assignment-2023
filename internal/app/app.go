package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app/config"
	v1 "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/controller/http/v1"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/storage"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segment"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/pkg/postgres"
	"github.com/gin-gonic/gin"
)

type onStopAction func(context.Context) error

type App struct {
	httpHandler   http.Handler
	config        *config.Config
	onStopActions []onStopAction
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	app := &App{
		config: cfg,
	}
	handler := gin.New()
	postgres, err := postgres.New(ctx, cfg.ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pg client: %w", err)
	}
	defer func() {
		if err != nil {
			postgres.Close(ctx)
		}
	}()
	app.onStopActions = append(app.onStopActions, func(ctx context.Context) error { return postgres.Close(ctx) })

	segmentStorage, err := storage.NewSegmentPsqlStorage(ctx, postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create pg client: %w", err)
	}
	segmentUseCase := segment.NewSegmentUsecase(segmentStorage)
	v1.NewRouter(handler, segmentUseCase)
	app.httpHandler = handler
	return app, nil
}

func (a App) Run(ctx context.Context) error {
	for _, act := range a.onStopActions {
		action := act
		defer func() {
			if err := action(ctx); err != nil {
				panic(err)
			}
		}()
	}
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", a.config.Port), a.httpHandler)
}
