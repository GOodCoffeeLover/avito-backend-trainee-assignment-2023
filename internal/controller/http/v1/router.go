package v1

import (
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segment"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, s segment.Segments) {
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
	{
		initSegmentRoutes(h, s)
	}
}
