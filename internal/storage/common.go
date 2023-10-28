package storage

import "errors"

var (
	ErrAlreadyExists = errors.New("alredy exists")
	ErrNotFound      = errors.New("not found")
	assigmentsTable  = struct {
		name        string
		userID      string
		segmentName string
		deleted     string
	}{
		name:        "assigments",
		userID:      "user_id",
		segmentName: "segment_name",
		deleted:     "deleted",
	}
)
