package types

import (
	"github.com/scottshotgg/reminder/pkg/reminder"
)

// DBReminder is used for the persistence layer
type (
	DBReminder struct {
		ID      string
		Created int64
		Moment  int64
		Message string
		Status  reminder.MsgStatus
		To      string
	}
)

// ToDB maps an API reminder to a DB reminder
func ToDB(r *reminder.Reminder) *DBReminder {
	return &DBReminder{
		ID:      r.ID,
		Created: r.Created,
		Moment:  r.Moment,
		Message: r.Message,
		Status:  r.Status,
		To:      r.To,
	}
}
