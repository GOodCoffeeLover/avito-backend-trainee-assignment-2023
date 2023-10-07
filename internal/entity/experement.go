package entity

import (
	"errors"
	"fmt"
)

type setOfSegmentNames map[SegmentName]struct{}

type Experement struct {
	UserID       UserID            `json:"user_id"`
	SegmentNames setOfSegmentNames `json:"segments_names"`
}

func NewExperiment(userID UserID, segments []SegmentName) (*Experement, error) {
	e := &Experement{UserID: userID}
	for _, segment := range segments {
		if err := e.AssignSegment(segment); errors.Is(err, ErrSegmentAlreadyAssigned) {
			return nil, ErrRepeatedSegmentForExperiment
		}
	}
	return e, nil
}

func (e *Experement) AssignSegment(segmentName SegmentName) error {
	if _, has := e.SegmentNames[segmentName]; has {
		return fmt.Errorf("%v already assigned to %v: %w", segmentName, e.UserID, ErrSegmentAlreadyAssigned)
	}
	e.SegmentNames[segmentName] = struct{}{}
	return nil
}
