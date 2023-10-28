package entity

import (
	"fmt"
	re "regexp"
)

var (
	segmentNameRegEx = re.MustCompile(segmentNameRegexPattern)
)

const (
	segmentNameMinimalLenght = 5
	segmentNameRegexPattern  = "^[A-Z0-9_]+$"
)

type SegmentName string

type Segment struct {
	Name SegmentName `json:"name"`
}

func NewSegment(name SegmentName) (*Segment, error) {

	if len(name) < segmentNameMinimalLenght {
		return nil, fmt.Errorf("%w: %v lenghts (%v) less than acceptable lenghts %v",
			ErrInvalidSegmentName, name, len(name), segmentNameMinimalLenght)
	}

	if !segmentNameRegEx.MatchString(string(name)) {
		return nil, fmt.Errorf("%w: %v not mathces pattern %v",
			ErrInvalidSegmentName, name, segmentNameRegEx)
	}

	return &Segment{
		Name: name,
	}, nil
}
