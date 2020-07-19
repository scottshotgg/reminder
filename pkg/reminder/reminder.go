package reminder

import (
	"time"
)

type Reminder interface {
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler

	GetTo() string
	GetCreated() int64
	// From() string
	// When() <some struct with all times>
	GetID() string
	GetMessage() string
	SetQueued(q bool)
	GetQueued() bool
	GetFired() bool
	GetUntil() time.Duration
	IsAfter(t time.Duration) bool
}
