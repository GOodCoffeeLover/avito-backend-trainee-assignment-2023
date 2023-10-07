package entity

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrDuplicatedValue = errors.New("duplicated value")
)

func (sos setOfSegmentNames) MarshalJSON() ([]byte, error) {
	sl := make([]SegmentName, 0, len(sos))
	for seg := range sos {
		sl = append(sl, seg)
	}
	return json.Marshal(sl)
}

func (sos setOfSegmentNames) UnmarshalJSON(bytes []byte) error {
	sl := []SegmentName{}
	if err := json.Unmarshal(bytes, &sl); err != nil {
		return err
	}
	for _, seg := range sl {
		if _, has := sos[seg]; has {
			return fmt.Errorf("%w: %v", ErrDuplicatedValue, seg)
		}
		sos[seg] = struct{}{}
	}
	return nil
}
