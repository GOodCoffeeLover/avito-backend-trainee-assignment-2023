package segment

import "github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"

type SegmentService interface {
	Create(*entity.Segment) error
	Delete(entity.SegmentName) error
}

type UserService interface {
	UnassingnSegmentAllUsers(entity.SegmentName) error
}
