package event

import (
	"context"
	"time"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type (
	EventStorage interface {
		ReadByUserID(context.Context, entity.UserID, *time.Time, *time.Time) ([]*entity.Event, error)
	}
	// TODO: maybe move to place, where it is used
	EventUseCase interface {
		ReadByUserID(context.Context, *InReadEventsByUserID) ([]*entity.Event, error)
	}
	UserStorage interface {
		ReadByID(context.Context, entity.UserID) (*entity.User, error)
	}
)
