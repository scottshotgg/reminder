package types

import (
	"time"

	"github.com/scottshotgg/reminder/pkg/reminder"
)

type DBReminder struct {
	ID      string
	Created int64
	Until   time.Duration
	Message string
	Queued  bool
	Fired   bool
	To      string
}

func ToDB(r reminder.Reminder) *DBReminder {
	return &DBReminder{
		ID:      r.GetID(),
		Created: r.GetCreated(),
		Until:   r.GetUntil(),
		Message: r.GetMessage(),
		Queued:  r.GetQueued(),
		Fired:   r.GetFired(),
		To:      r.GetTo(),
	}
}

type internalV1 struct {
	ID      string
	Created int64
	Until   time.Duration
}
