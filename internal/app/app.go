package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app/config"
	http_router "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/controller/http"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/storage"
	assignment_usecase "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/assignment"
	segment_usecase "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segmnet"
	user_usecase "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/user"
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
	var err error
	defer func() {
		if err != nil {
			app.stop(ctx)
		}
	}()

	postgres, trm, err := postgres.New(ctx, cfg.ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pg client: %w", err)
	}
	app.onStopActions = append(app.onStopActions, func(ctx context.Context) error { postgres.Close(ctx); return nil })

	segmentStorage, err := storage.NewSegmentPsql(ctx, postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create pg segments storage: %w", err)
	}
	segments := segment_usecase.New(segmentStorage, trm)

	userStorage, err := storage.NewUserPsql(ctx, postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create pg users storage: %w", err)
	}
	users := user_usecase.New(userStorage, trm)

	assignmentStorage, err := storage.NewAssignmentPsql(ctx, postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to create pg assignment storage: %w", err)
	}
	assignments := assignment_usecase.New(segmentStorage, userStorage, assignmentStorage, trm)

	handler := gin.New()
	http_router.NewRouter(handler, segments, users, assignments)
	app.httpHandler = handler
	return app, nil
}

func (a App) Run(ctx context.Context) error {
	defer func() {
		a.stop(ctx)
	}()
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", a.config.Port), a.httpHandler)
}

func (a App) stop(ctx context.Context) {
	for _, action := range a.onStopActions {
		if err := action(ctx); err != nil {
			panic(err)
		}
	}
}
