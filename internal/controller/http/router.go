package http

import (
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/assignment"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/event"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segment"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/user"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter(handler *gin.Engine, segments segment.SegmentUseCase, users user.UserUseCase, asassignments assignment.AssigmentUseCase, events event.EventUseCase) {
	// Options
	handler.Use(gin.Logger(), gin.Recovery())

	// Swagger
	// swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	// handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/ready", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := handler.Group("/v1")

	initSegmentRoutes(h, segments)
	initUserRoutes(h, users)
	initAssigmentRoutes(h, asassignments)
	initEventRoutes(h, events)

}
