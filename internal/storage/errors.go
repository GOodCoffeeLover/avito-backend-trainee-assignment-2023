package storage

import "errors"

var (
	ErrAlreadyExists = errors.New("alredy exists")
	ErrNotFound      = errors.New("not found")
)
