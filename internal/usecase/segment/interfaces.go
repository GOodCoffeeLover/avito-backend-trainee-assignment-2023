package segment

import (
	"context"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type (
	SegmentStorage interface {
		Create(context.Context, *entity.Segment) error
		Read(context.Context, entity.SegmentName) (*entity.Segment, error)
		Delete(context.Context, entity.SegmentName) error
	}
	ExperementStorage interface {
		UnassingnSegmentFromAllUsers(context.Context, entity.SegmentName) error
	}
	Segments interface {
		Create(context.Context, entity.SegmentName) error
		Read(context.Context, entity.SegmentName) (*entity.Segment, error)
		Delete(context.Context, entity.SegmentName) error
	}
)
