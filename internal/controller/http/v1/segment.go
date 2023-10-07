package v1

import (
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/segment"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
}

const (
	segmentNameParam = "segment_name"
)

type segmentRoutes struct {
	segmets segment.Segments
}

func initSegmentRoutes(handler *gin.RouterGroup, segments segment.Segments) {
	sr := segmentRoutes{segmets: segments}
	h := handler.Group("/segments")
	h.GET("/:"+segmentNameParam, sr.read)
	h.POST("/", sr.create)
	h.DELETE("/:"+segmentNameParam, sr.delete)
}

func (sr *segmentRoutes) read(ctx *gin.Context) {
	name := ctx.Param(segmentNameParam)

	segment, err := sr.segmets.Read(ctx.Request.Context(), entity.SegmentName(name))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, segment)

}
func (sr *segmentRoutes) create(ctx *gin.Context) {
	seg := entity.Segment{}
	if err := ctx.ShouldBindJSON(&seg); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, errorResponse{err.Error()})
		return
	}

	if err := sr.segmets.Create(ctx.Request.Context(), seg.Name); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func (sr *segmentRoutes) delete(ctx *gin.Context) {
	name := ctx.Param(segmentNameParam)

	err := sr.segmets.Delete(ctx.Request.Context(), entity.SegmentName(name))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{err.Error()})
		return
	}
	ctx.Status(http.StatusOK)

}
