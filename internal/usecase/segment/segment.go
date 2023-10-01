package segment

import (
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type SegmentUseCase struct {
	segmentService SegmentService
	userService    UserService
}

func New(segmentService SegmentService, userService UserService) *SegmentUseCase {
	return &SegmentUseCase{
		segmentService: segmentService,
		userService:    userService,
	}
}

func (suc *SegmentUseCase) Create(segment *entity.Segment) error {
	err := suc.segmentService.Create(segment)
	if err != nil {
		return fmt.Errorf("failed to create segment(%v): %w", segment.Name(), err)
	}
	return nil
}

func (suc *SegmentUseCase) Delete(segmentName entity.SegmentName) error {
	err := suc.userService.UnassingnSegmentAllUsers(segmentName)
	if err != nil {
		return fmt.Errorf("failed to unassgn segment %v from users: %w", segmentName, err)
	}

	// TODO: fix problem when all users unassigned, but segment not deleted
	err = suc.segmentService.Delete(segmentName)
	if err != nil {
		return fmt.Errorf("failed to delete segment(%v): %w", segmentName, err)
	}
	return nil
}
