package entity

import (
	"errors"
	"fmt"
)

type Experiment struct {
	UserID       UserID        `json:"user_id"`
	SegmentNames []SegmentName `json:"segment_names"`
}

func NewExperiment(userID UserID, segmentNames []SegmentName) (*Experiment, error) {
	e := &Experiment{UserID: userID}
	for _, segment := range segmentNames {
		if err := e.AssignSegment(segment); errors.Is(err, ErrSegmentAlreadyAssigned) {
			return nil, ErrRepeatedSegment
		}
	}
	return e, nil
}

func (u *Experiment) AssignSegment(segmentName SegmentName) error {
	for _, seg := range u.SegmentNames {
		if seg == segmentName {
			return fmt.Errorf("%v already assigned to %v: %w", segmentName, u.UserID, ErrSegmentAlreadyAssigned)
		}
	}
	u.SegmentNames = append(u.SegmentNames, segmentName)
	return nil
}
