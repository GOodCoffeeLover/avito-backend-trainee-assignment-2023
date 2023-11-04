package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	users user.UserUseCase
}

func initUserRoutes(handler *gin.RouterGroup, users user.UserUseCase) {
	ur := newUserRoutes(users)
	h := handler.Group("/user")

	h.GET("/", ur.readAll)
	h.GET(fmt.Sprintf("/:%v", userIDParam), ur.read)

	h.POST(fmt.Sprintf("/:%v", userIDParam), ur.create)
	h.DELETE(fmt.Sprintf("/:%v", userIDParam), ur.delete)
}

func newUserRoutes(users user.UserUseCase) userRoutes {
	return userRoutes{users: users}
}
func (ur *userRoutes) read(ctx *gin.Context) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", entity.ErrInvalidArgument, err))
		return
	}
	segment, err := ur.users.Read(ctx.Request.Context(), entity.UserID(uid))
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, segment)

}

func (ur *userRoutes) readAll(ctx *gin.Context) {

	users, err := ur.users.ReadAll(ctx.Request.Context())
	if err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, users)

}

func (ur *userRoutes) create(ctx *gin.Context) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", entity.ErrInvalidArgument, err))
		return
	}
	if err := ur.users.Create(ctx.Request.Context(), entity.UserID(uid)); err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (ur *userRoutes) delete(ctx *gin.Context) {
	uid, err := strconv.ParseUint(ctx.Param(userIDParam), 10, 32)
	if err != nil {
		abortWithErrorAnalize(ctx, fmt.Errorf("%w: %w", entity.ErrInvalidArgument, err))
		return
	}
	if err := ur.users.Delete(ctx.Request.Context(), entity.UserID(uid)); err != nil {
		abortWithErrorAnalize(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)

}
