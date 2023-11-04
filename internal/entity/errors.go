package entity

import "errors"

var (
	ErrInvalidSegmentName     = errors.New("ivalid segment name")
	ErrSegmentAlreadyAssigned = errors.New("segment already assigned")
	ErrRepeatedSegment        = errors.New("repeated segment")
	ErrDuplicatedValue        = errors.New("duplicated value")
	ErrNotFound               = errors.New("not found")
	ErrAlreadyExists          = errors.New("alredy exists")
)
