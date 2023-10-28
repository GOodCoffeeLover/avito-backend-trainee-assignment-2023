package v1

import (
	"errors"
	"net/http"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/storage"
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

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

func abortWithErrorAnalize(ctx *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	// case errors.Is(err, pgx.ErrNoRows):
	//     status = http.StatusNotFound
	case errors.Is(err, storage.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, storage.ErrAlreadyExists):
		status = http.StatusBadRequest
	case errors.Is(err, ErrInvalidArgument):
		status = http.StatusBadRequest
	}
	ctx.AbortWithStatusJSON(status, errorResponse{err.Error(), status})
}
