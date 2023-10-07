package entity

import (
	"fmt"
	re "regexp"
)

var (
	segmentNameRegEx *re.Regexp = nil
)

const (
	segmentNameMinimalLenght = 5
	segmentNameRegexPattern  = "[A-Z0-9_]+"
)

func init() {
	var err error
	segmentNameRegEx, err = re.Compile(segmentNameRegexPattern)
	if err != nil {
		panic(err)
	}
}

type SegmentName string

type Segment struct {
	Name SegmentName
}

func NewSegment(name SegmentName) (*Segment, error) {

	if len(name) < segmentNameMinimalLenght {
		return nil, fmt.Errorf("%w: %v lenghts (%v) less than acceptable lenghts %v",
			ErrInvalidSegmentNameLenght, name, len(name), segmentNameMinimalLenght)
	}

	if !segmentNameRegEx.MatchString(string(name)) {
		return nil, fmt.Errorf("%w: %v not mathces %v",
			ErrInvalidSegmentNamePattern, name, segmentNameRegEx)
	}

	return &Segment{
		Name: name,
	}, nil
}
