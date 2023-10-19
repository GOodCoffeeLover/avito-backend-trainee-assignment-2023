package v1

import (
	"net/http"

	segment "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segmnet"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, segments segment.SegmentUseCase, users user.UserUseCase) {
	// Options
	handler.Use(gin.Logger(), gin.Recovery())

	// Swagger
	// swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	// handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")

	initSegmentRoutes(h, segments)
	initUserRoutes(h, users)

}
