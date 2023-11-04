package assignment

import (
	"context"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type (
	SegmentStorage interface {
		ReadByName(context.Context, entity.SegmentName) (*entity.Segment, error)
	}
	UserStorage interface {
		ReadByID(context.Context, entity.UserID) (*entity.User, error)
	}
	AssignmentStorage interface {
		ReadByUserID(context.Context, entity.UserID) ([]*entity.Assignment, error)
		Create(context.Context, *entity.Assignment) error
		Delete(context.Context, *entity.Assignment) error
	}
	EventsStorage interface {
		Create(context.Context, *entity.Event) error
	}
	// TODO: maybe move to place, where it is used
	AssigmentUseCase interface {
		ReadByUserID(context.Context, entity.UserID) ([]*entity.Assignment, error)
		SetToUserByID(context.Context, entity.UserID, []entity.SegmentName) error
		UnsetToUserByID(context.Context, entity.UserID, []entity.SegmentName) error
	}
)
