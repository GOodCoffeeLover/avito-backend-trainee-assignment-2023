package segment

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type SegmentUseCase struct {
	segmentStorage SegmentStorage
	userStorage    UserStorage
}

func New(segmentStorage SegmentStorage, userStorage UserStorage) *SegmentUseCase {
	return &SegmentUseCase{
		segmentStorage: segmentStorage,
		userStorage:    userStorage,
	}
}

func (suc *SegmentUseCase) Create(ctx context.Context, segmentName entity.SegmentName) error {
	segment, err := entity.NewSegment(segmentName)
	if err != nil {
		return fmt.Errorf("failed to build new segment: %w", err)
	}

	err = suc.segmentStorage.Create(ctx, segment)
	if err != nil {
		return fmt.Errorf("failed to create segment(%v): %w", segmentName, err)
	}
	return nil
}

func (suc *SegmentUseCase) Delete(ctx context.Context, segmentName entity.SegmentName) error {
	err := suc.segmentStorage.Delete(ctx, segmentName)
	if err != nil {
		return fmt.Errorf("failed to delete segment(%v): %w", segmentName, err)
	}

	// TODO: fix problem when all users unassigned, but segment not deleted
	err = suc.userStorage.UnassingnSegmentAllUsers(ctx, segmentName)
	if err != nil {
		return fmt.Errorf("failed to unassgn segment %v from users: %w", segmentName, err)
	}
	return nil
}
