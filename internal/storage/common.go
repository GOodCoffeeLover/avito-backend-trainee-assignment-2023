package storage

import "errors"

var (
	ErrAlreadyExists = errors.New("alredy exists")
	ErrNotFound      = errors.New("not found")
	assignmentsTable = struct {
		name        string
		userID      string
		segmentName string
		deleted     string
	}{
		name:        "assignments",
		userID:      "user_id",
		segmentName: "segment_name",
		deleted:     "deleted",
	}
)
