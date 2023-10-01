package entity

import (
	"errors"
	"fmt"
	re "regexp"
)

var (
	segmentNameRegEx             *re.Regexp = nil
	ErrInvalidSegmentNamePattern            = errors.New("ivalid segment pattern")
	ErrInvalidSegmentNameLenght             = errors.New("ivalid segment lenght")
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
	name SegmentName
}

func NewSegment(name SegmentName) (*Segment, error) {

	if len(name) < segmentNameMinimalLenght {
		return nil,
			fmt.Errorf("%w: %v lenghts (%v) less than acceptable lenghts %v",
				ErrInvalidSegmentNameLenght, name, len(name), segmentNameMinimalLenght)
	}

	if !segmentNameRegEx.MatchString(string(name)) {
		return nil, fmt.Errorf("%w: %v not mathces %v",
			ErrInvalidSegmentNamePattern, name, segmentNameRegEx)
	}

	return &Segment{
		name: name,
	}, nil
}

func (s Segment) Name() SegmentName {
	return s.name
}
