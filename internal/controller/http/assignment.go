package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	assignment "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/assignment"
	"github.com/gin-gonic/gin"
)

type assignmentRoutes struct {
	assignments assignment.AssigmentUseCase
}

func initAssigmentRoutes(handler *gin.RouterGroup, assignments assignment.AssigmentUseCase) {
	ar := newAssigmentRoutes(assignments)
	h := handler.Group("/user")

	h.GET(fmt.Sprintf("/:%v/assignments", userIDParam), ar.getByUser)
	// not very REST way to set list
	h.POST(fmt.Sprintf("/:%v/assignments", userIDParam), ar.assignListSegmentsToUser)
	h.DELETE(fmt.Sprintf("/:%v/assignments", userIDParam), ar.unassignListSegmentsToUser)
}

func newAssigmentRoutes(assignments assignment.AssigmentUseCase) assignmentRoutes {
	return assignmentRoutes{assignments: assignments}
}

func (ar *assignmentRoutes) getByUser(ctx *gin.Context) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
		return
	}
	assignments, err := ar.assignments.ReadByUserID(ctx.Request.Context(), entity.UserID(uid))
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, assignments)

}

func (ar *assignmentRoutes) assignListSegmentsToUser(ctx *gin.Context) {
	uid, err := ar.retriveUserID(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
		return
	}
	segs, err := ar.retriveSegmentsNames(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
	}

	err = ar.assignments.SetToUserByID(ctx.Request.Context(), uid, segs)
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (ar *assignmentRoutes) unassignListSegmentsToUser(ctx *gin.Context) {
	uid, err := ar.retriveUserID(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
		return
	}
	segs, err := ar.retriveSegmentsNames(ctx)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", ErrInvalidArgument, err))
	}

	err = ar.assignments.UnsetToUserByID(ctx.Request.Context(), uid, segs)
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (ar *assignmentRoutes) retriveUserID(ctx *gin.Context) (entity.UserID, error) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	return entity.UserID(uid), err
}

func (ar *assignmentRoutes) retriveSegmentsNames(ctx *gin.Context) ([]entity.SegmentName, error) {
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
