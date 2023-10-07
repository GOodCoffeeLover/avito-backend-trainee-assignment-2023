package storage

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
)

type SegmentInMemoryStorage struct {
	segments map[entity.SegmentName]entity.Segment
}

func NewSegmentInMemoryStorage() *SegmentInMemoryStorage {
	return &SegmentInMemoryStorage{
		segments: make(map[entity.SegmentName]entity.Segment),
	}
}

func (sims SegmentInMemoryStorage) Create(ctx context.Context, segment *entity.Segment) error {
	if segment == nil {
		return fmt.Errorf("nil segment")
	}

	if _, ok := sims.segments[segment.Name]; ok {
		return fmt.Errorf("already exists")
	}
	sims.segments[segment.Name] = *segment
	return nil
}

func (sims SegmentInMemoryStorage) Read(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {

	segment, ok := sims.segments[segmentName]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return &segment, nil
}

func (sims SegmentInMemoryStorage) Delete(ctx context.Context, segmentName entity.SegmentName) error {

	if _, ok := sims.segments[segmentName]; !ok {
		return fmt.Errorf("not found")
	}
	delete(sims.segments, segmentName)
	return nil
}
