package entity

type UserID uint

type User struct {
	id       UserID
	segments map[SegmentName]Segment
}

func NewUser(id UserID) (*User, error) {
	return &User{
		id:       id,
		segments: make(map[SegmentName]Segment),
	}, nil
}

func (u *User) AssignSegments(segments ...Segment) {
	for _, segment := range segments {
		s := segment
		u.segments[s.name] = s
	}
}

func (u *User) ListAssignedSegments() []Segment {
	segments := []Segment{}
	for _, segment := range u.segments {
		segments = append(segments, segment)
	}
	return segments
}

func (u *User) ID() UserID {
	return u.id
}
