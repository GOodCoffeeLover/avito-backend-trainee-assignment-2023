package assigment

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
		ReadByUserID(context.Context, entity.UserID) (*entity.Assignment, error)
		Save(context.Context, *entity.Assignment) error
	}
	// TODO: maybe move to place, where it is used
	AssigmentUseCase interface {
		ReadByUserID(context.Context, entity.UserID) (*entity.Assignment, error)
		AssignSegments(context.Context, entity.UserID, []entity.SegmentName) error
		UnassignSegments(context.Context, entity.UserID, []entity.SegmentName) error
	}
)
