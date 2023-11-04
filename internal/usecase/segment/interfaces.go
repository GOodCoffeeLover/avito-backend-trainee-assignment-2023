package segment

import (
	"context"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type (
	SegmentStorage interface {
		Create(context.Context, *entity.Segment) error
		ReadByName(context.Context, entity.SegmentName) (*entity.Segment, error)
		ReadAll(context.Context) ([]*entity.Segment, error)
		Delete(context.Context, entity.SegmentName) error
	}
	EventStorage interface {
		Create(context.Context, *entity.Event) error
	}
	// TODO: maybe move to place, where it is used
	SegmentUseCase interface {
		Create(context.Context, entity.SegmentName) error
		Read(context.Context, entity.SegmentName) (*entity.Segment, error)
		ReadAll(context.Context) ([]*entity.Segment, error)
		Delete(context.Context, entity.SegmentName) error
	}
)
