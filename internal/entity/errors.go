package entity

import "errors"

var (
	ErrInvalidSegmentName = errors.New("ivalid segment name")
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("alredy exists")
	ErrInvalidArgument    = errors.New("invalid argument")
)
