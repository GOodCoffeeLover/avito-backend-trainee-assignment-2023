package entity

type UserID uint

type User struct {
	id       UserID
	segments []SegmentName
}

func NewUser(id UserID) (*User, error) {
	return &User{
		id:       id,
		segments: make([]SegmentName, 0),
	}, nil
}

func (u *User) AssignSegments(segments ...SegmentName) {
	u.segments = append(u.segments, segments...)
}

func (u *User) ListAssignedSegments() []SegmentName {
	res := make([]SegmentName, len(u.segments))
	copy(res, u.segments)
	return res
}

func (u *User) ID() UserID {
	return u.id
}
