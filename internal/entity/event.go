package entity

import (
	"time"
)

type Operation string

const (
	Created Operation = "created"
	Deleted Operation = "deleted"
)

type Event struct {
	User    *UserID
	Segment *SegmentName
	Op      Operation
	Ts      time.Time
}

func NewEvent(uid *UserID, segment *SegmentName, operation Operation) *Event {
	return &Event{
		User:    uid,
		Segment: segment,
		Op:      operation,
		Ts:      time.Now(),
	}
}

// func (e Event) String() string {
// 	uid := "all"
// 	if e.uid != nil {
// 		uid = fmt.Sprintf("%v", *e.uid)
// 	}
// 	segment := "all"
// 	if e.segment != nil {
// 		segment = fmt.Sprintf("%v", *e.segment)
// 	}
// 	return fmt.Sprintf("%v;%v;%v;%v", uid, segment, e.operation, e.ts)
// }
