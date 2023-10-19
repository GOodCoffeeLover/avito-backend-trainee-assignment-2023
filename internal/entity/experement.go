package entity

import (
	"errors"
	"fmt"
)

type UserID uint

type User struct {
	UserID       UserID        `json:"user_id"`
	SegmentNames []SegmentName `json:"segment_names"`
}

func NewExperiment(userID UserID, segmentNames []SegmentName) (*User, error) {
	u := &User{UserID: userID}
	for _, segment := range segmentNames {
		if err := u.AssignSegment(segment); errors.Is(err, ErrSegmentAlreadyAssigned) {
			return nil, ErrRepeatedSegmentForExperiment
		}
	}
	return u, nil
}

func (u *User) AssignSegment(segmentName SegmentName) error {
	for _, seg := range u.SegmentNames {
		if seg == segmentName {
			return fmt.Errorf("%v already assigned to %v: %w", segmentName, u.UserID, ErrSegmentAlreadyAssigned)
		}
	}
	u.SegmentNames = append(u.SegmentNames, segmentName)
	return nil
}
