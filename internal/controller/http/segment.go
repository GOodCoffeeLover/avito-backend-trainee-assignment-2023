package http

import (
	"fmt"
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segment"
	"github.com/gin-gonic/gin"
)

type segmentRoutes struct {
	segmets segment.SegmentUseCase
}

func initSegmentRoutes(handler *gin.RouterGroup, segments segment.SegmentUseCase) {
	sr := newSegmentRoutes(segments)
	h := handler.Group("/segment")

	h.GET("/", sr.readAll)
	h.GET(fmt.Sprintf("/:%v", segmentNameParam), sr.read)

	h.POST(fmt.Sprintf("/:%v", segmentNameParam), sr.create)
	h.DELETE(fmt.Sprintf("/:%v", segmentNameParam), sr.delete)
}

func newSegmentRoutes(segments segment.SegmentUseCase) segmentRoutes {
	return segmentRoutes{segmets: segments}
}
func (sr *segmentRoutes) read(ctx *gin.Context) {
	name := ctx.Param(segmentNameParam)

	segment, err := sr.segmets.Read(ctx.Request.Context(), entity.SegmentName(name))
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, segment)

}

func (sr *segmentRoutes) readAll(ctx *gin.Context) {

	segments, err := sr.segmets.ReadAll(ctx.Request.Context())
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, segments)

}

func (sr *segmentRoutes) create(ctx *gin.Context) {
	name := ctx.Param(segmentNameParam)

	if err := sr.segmets.Create(ctx.Request.Context(), entity.SegmentName(name)); err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (sr *segmentRoutes) delete(ctx *gin.Context) {
	name := ctx.Param(segmentNameParam)

	if err := sr.segmets.Delete(ctx.Request.Context(), entity.SegmentName(name)); err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)

}
