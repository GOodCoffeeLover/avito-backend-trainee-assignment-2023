package segment

import (
	"context"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type SegmentStorage interface {
	Create(context.Context, *entity.Segment) error
	Delete(context.Context, entity.SegmentName) error
}

type UserStorage interface {
	UnassingnSegmentAllUsers(context.Context, entity.SegmentName) error
}
