package storage

var (
	assignments = struct {
		table       string
		userID      string
		segmentName string
		deleted     string
	}{
		table:       "assignments",
		userID:      "user_id",
		segmentName: "segment_name",
		deleted:     "deleted",
	}
	segments = struct {
		table   string
		name    string
		deleted string
	}{
		table:   "segments",
		name:    "name",
		deleted: "deleted",
	}
	users = struct {
		table   string
		id      string
		deleted string
	}{
		table:   "users",
		id:      "id",
		deleted: "deleted",
	}
)
