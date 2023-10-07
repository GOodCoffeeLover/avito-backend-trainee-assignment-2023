package app

import (
	"fmt"
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app/config"
	v1 "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/controller/http/v1"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/storage"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segment"
	"github.com/gin-gonic/gin"
)

type App struct {
	httpHandler http.Handler
	config      *config.Config
}

func New(cfg *config.Config) (*App, error) {
	handler := gin.New()
	segmentStorage := storage.NewSegmentInMemoryStorage()
	segmentUseCase := segment.NewSegmentUsecase(segmentStorage)
	v1.NewRouter(handler, segmentUseCase)
	return &App{
		httpHandler: handler,
		config:      cfg,
	}, nil
}

func (a App) Run() error {
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", a.config.Port), a.httpHandler)
}
