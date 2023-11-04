package segment

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

type UseCase struct {
	segments   SegmentStorage
	events     EventStorage
	trxManager *manager.Manager
}

func New(segments SegmentStorage, events EventStorage, trxManager *manager.Manager) *UseCase {
	return &UseCase{
		segments:   segments,
		events:     events,
		trxManager: trxManager,
	}
}

func (uc *UseCase) Create(ctx context.Context, segmentName entity.SegmentName) error {
	segment, err := entity.NewSegment(segmentName)
	if err != nil {
		return fmt.Errorf("failed to build segment: %w", err)
	}
	err = uc.trxManager.Do(ctx, func(ctx context.Context) error {
		return uc.segments.Create(ctx, segment)
	})
	if err != nil {
		return fmt.Errorf("failed to create segment %v: %w", segmentName, err)
	}
	return nil
}

func (uc *UseCase) Read(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {
	var seg *entity.Segment
	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		seg, err = uc.segments.ReadByName(ctx, segmentName)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read segment %v: %w", segmentName, err)
	}
	return seg, nil
}

func (uc *UseCase) ReadAll(ctx context.Context) ([]*entity.Segment, error) {
	var segs []*entity.Segment
	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		segs, err = uc.segments.ReadAll(ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read all segments: %w", err)
	}
	return segs, nil
}

func (uc *UseCase) Delete(ctx context.Context, segmentName entity.SegmentName) error {
	return uc.trxManager.Do(ctx, func(ctx context.Context) error {

		if err := uc.segments.Delete(ctx, segmentName); err != nil {
			return fmt.Errorf("failed to delete segment %v: %w", segmentName, err)
		}

		if err := uc.events.Create(ctx, entity.NewEvent(nil, &segmentName, entity.Deleted)); err != nil {
			return fmt.Errorf("failed to save event of segment %v deletion: %w", segmentName, err)
		}
		return nil
	})
}
