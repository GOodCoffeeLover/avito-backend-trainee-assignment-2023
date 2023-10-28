package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	assigment "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/assigment"
	"github.com/gin-gonic/gin"
)

type assigmentRoutes struct {
	assigments assigment.AssigmentUseCase
}

func initAssigmentRoutes(handler *gin.RouterGroup, assigments assigment.AssigmentUseCase) {
	ar := newAssigmentRoutes(assigments)
	h := handler.Group("/user")

	h.GET(fmt.Sprintf("/:%v/assigments", userIDParam), ar.getByUser)
	// not very REST way to set list
	h.POST(fmt.Sprintf("/:%v/assigments", userIDParam), ar.assignListSegmentsToUser)
	h.DELETE(fmt.Sprintf("/:%v/assigments", userIDParam), ar.unassignListSegmentsToUser)
}

func newAssigmentRoutes(assigments assigment.AssigmentUseCase) assigmentRoutes {
	return assigmentRoutes{assigments: assigments}
}

func (ar *assigmentRoutes) getByUser(ctx *gin.Context) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
		return
	}
	assigments, err := ar.assigments.ReadByUserID(ctx.Request.Context(), entity.UserID(uid))
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, assigments)

}

func (ar *assigmentRoutes) assignListSegmentsToUser(ctx *gin.Context) {
	uid, err := ar.retriveUserID(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
		return
	}
	segs, err := ar.retriveSegmentsNames(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
	}

	err = ar.assigments.SetToUserByID(ctx.Request.Context(), uid, segs)
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (ar *assigmentRoutes) unassignListSegmentsToUser(ctx *gin.Context) {
	uid, err := ar.retriveUserID(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
		return
	}
	segs, err := ar.retriveSegmentsNames(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
	}

	err = ar.assigments.UnsetToUserByID(ctx.Request.Context(), uid, segs)
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (ar *assigmentRoutes) retriveUserID(ctx *gin.Context) (entity.UserID, error) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	return entity.UserID(uid), err
}

func (ar *assigmentRoutes) retriveSegmentsNames(ctx *gin.Context) ([]entity.SegmentName, error) {
	rawBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read body: %w", err)
	}
	segs := struct {
		Segments []entity.SegmentName `json:"segments"`
	}{}
	err = json.Unmarshal(rawBody, &segs)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal segments: %w", err)
	}
	return segs.Segments, nil
}
