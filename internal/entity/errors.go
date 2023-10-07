package entity

import "errors"

var (
	ErrInvalidSegmentNamePattern    = errors.New("ivalid segment pattern")
	ErrInvalidSegmentNameLenght     = errors.New("ivalid segment lenght")
	ErrSegmentAlreadyAssigned       = errors.New("segment already assigned")
	ErrRepeatedSegmentForExperiment = errors.New("repeated segment for experiment")
)
