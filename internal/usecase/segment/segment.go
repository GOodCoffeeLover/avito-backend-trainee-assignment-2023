package segment

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type SegmentUseCase struct {
	segmentStorage SegmentStorage
	// experementStorage ExperementStorage
}

func NewSegmentUsecase(segmentStorage SegmentStorage /*, experementStorage ExperementStorage*/) *SegmentUseCase {
	return &SegmentUseCase{
		segmentStorage: segmentStorage,
		// experementStorage: experementStorage,
	}
}

func (suc *SegmentUseCase) Create(ctx context.Context, segmentName entity.SegmentName) error {
	segment, err := entity.NewSegment(segmentName)
	if err != nil {
		return fmt.Errorf("failed to build new segment: %w", err)
	}

	err = suc.segmentStorage.Create(ctx, segment)
	if err != nil {
		return fmt.Errorf("failed to create %v: %w", segmentName, err)
	}
	return nil
}

func (suc *SegmentUseCase) Read(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {
	seg, err := suc.segmentStorage.Read(ctx, segmentName)
	if err != nil {
		return nil, fmt.Errorf("failed to read %v: %w", segmentName, err)
	}
	return seg, nil
}

func (suc *SegmentUseCase) Delete(ctx context.Context, segmentName entity.SegmentName) error {
	err := suc.segmentStorage.Delete(ctx, segmentName)
	if err != nil {
		return fmt.Errorf("failed to delete %v: %w", segmentName, err)
	}

	// TODO: fix problem when all users unassigned, but segment not deleted
	// err = suc.experementStorage.UnassingnSegmentFromAllUsers(ctx, segmentName)
	// if err != nil {
	// 	return fmt.Errorf("failed to unassgn segment %v from users: %w", segmentName, err)
	// }
	return nil
}
