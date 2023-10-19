package segment

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

type SegmentUseCase struct {
	segmentStorage SegmentStorage
	trxManager     *manager.Manager // for history transactions
}

func NewSegmentUsecase(segmentStorage SegmentStorage, trxManager *manager.Manager) *SegmentUseCase {
	return &SegmentUseCase{
		segmentStorage: segmentStorage,
		trxManager:     trxManager,
	}
}

func (suc *SegmentUseCase) Create(ctx context.Context, segmentName entity.SegmentName) error {
	segment, err := entity.NewSegment(segmentName)
	if err != nil {
		return fmt.Errorf("failed to build new segment: %w", err)
	}
	err = suc.trxManager.Do(ctx, func(ctx context.Context) error {
		return suc.segmentStorage.Create(ctx, segment)
	})
	if err != nil {
		return fmt.Errorf("failed to create %v: %w", segmentName, err)
	}
	return nil
}

func (suc *SegmentUseCase) Read(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {
	var seg *entity.Segment
	err := suc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		seg, err = suc.segmentStorage.ReadByName(ctx, segmentName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read %v: %w", segmentName, err)
	}
	return seg, nil
}

func (suc *SegmentUseCase) ReadAll(ctx context.Context) ([]*entity.Segment, error) {
	var segs []*entity.Segment
	err := suc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		segs, err = suc.segmentStorage.ReadAll(ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read all segments: %w", err)
	}
	return segs, nil
}

func (suc *SegmentUseCase) Delete(ctx context.Context, segmentName entity.SegmentName) error {
	err := suc.segmentStorage.Delete(ctx, segmentName)
	if err != nil {
		return fmt.Errorf("failed to delete %v: %w", segmentName, err)
	}
	return nil
}
