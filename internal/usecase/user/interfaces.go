package user

import (
	"context"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type (
	UserStorage interface {
		Create(context.Context, *entity.User) error
		ReadByID(context.Context, entity.UserID) (*entity.User, error)
		ReadAll(context.Context) ([]*entity.User, error)
		Delete(context.Context, entity.UserID) error
	}
	UserUseCase interface {
		Create(context.Context, entity.UserID) error
		Read(context.Context, entity.UserID) (*entity.User, error)
		ReadAll(context.Context) ([]*entity.User, error)
		Delete(context.Context, entity.UserID) error
	}
)
