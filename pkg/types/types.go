package types

import (
	"github.com/scottshotgg/reminder/pkg/reminder"
)

type DBReminder struct {
	ID      string
	Created int64
	Moment  int64
	Message string
	Status  reminder.MsgStatus
	To      string
}

func ToDB(r reminder.Reminder) *DBReminder {
	return &DBReminder{
		ID:      r.GetID(),
		Created: r.GetCreated(),
		Moment:  r.GetMoment(),
		Message: r.GetMessage(),
		Status:  r.GetStatus(),
		To:      r.GetTo(),
	}
}
