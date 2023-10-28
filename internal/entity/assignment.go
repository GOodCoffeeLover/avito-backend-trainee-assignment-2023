package entity

type Assignment struct {
	User    UserID      `json:"user_id"`
	Segment SegmentName `json:"segment_names"`
}

func NewAssignment(user UserID, segment SegmentName) (*Assignment, error) {
	return &Assignment{
		User:    user,
		Segment: segment,
	}, nil
}
