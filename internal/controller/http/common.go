package http

import (
	"errors"
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

const (
	segmentNameParam = "segment_name"
	userIDParam      = "user_id"
)

func abortWithErrorAnalize(ctx *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, entity.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, entity.ErrAlreadyExists):
		status = http.StatusBadRequest
	case errors.Is(err, entity.ErrInvalidArgument):
		status = http.StatusBadRequest
	}
	ctx.AbortWithStatusJSON(status, errorResponse{err.Error(), status})
}
