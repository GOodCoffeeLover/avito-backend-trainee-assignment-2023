package entity

import (
	"errors"
	"fmt"
)

type Assignment struct {
	UserID       UserID        `json:"user_id"`
	SegmentNames []SegmentName `json:"segment_names"`
}

func NewAssignment(userID UserID, segmentNames []SegmentName) (*Assignment, error) {
	a := &Assignment{
		UserID:       userID,
		SegmentNames: make([]SegmentName, 0, len(segmentNames)),
	}
	for _, segment := range segmentNames {
		if err := a.AssignSegment(segment); errors.Is(err, ErrSegmentAlreadyAssigned) {
			return nil, ErrRepeatedSegment
		}
	}
	return a, nil
}

func (a *Assignment) AssignSegment(segmentName SegmentName) error {
	for _, seg := range a.SegmentNames {
		if seg == segmentName {
			return fmt.Errorf("%v already assigned to %v: %w", segmentName, a.UserID, ErrSegmentAlreadyAssigned)
		}
	}
	a.SegmentNames = append(a.SegmentNames, segmentName)
	return nil
}

func (a *Assignment) UnssignSegment(segName SegmentName) error {
	for i, seg := range a.SegmentNames {
		if seg == segName {
			a.SegmentNames = append(a.SegmentNames[0:i], a.SegmentNames[i+1:len(a.SegmentNames)]...)
			return nil
		}
	}
	return fmt.Errorf("%w segment %v", ErrNotFound, segName)
}
