package v1

import (
	"time"
)

type V1 struct {
	Created int64
	After   time.Duration
	Msg     string
	Fired   bool
	Queued  bool
	To      string
}

// func NewV1(created int64, after time.Duration, message string, to string) *V1 {
// 	return &V1{
// 		created: created,
// 		after:   after,
// 		message: message,
// 		to:      to,
// 	}
// }

func (v *V1) Message() string {
	return v.Msg
}

// func (v V1) MarshalBinary() ([]byte, error) {
// 	return binary.Marshal(v)
// }

// func (v V1) UnmarshalBinary(data []byte) error {
// 	return binary.Unmarshal(data, v)
// }
