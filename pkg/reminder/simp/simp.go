package simp

import (
	"github.com/scottshotgg/reminder/pkg/reminder"
)

type SimpleReminder struct {
	message string
}

func (r *SimpleReminder) Message() string {
	return r.message
}

func New(msg string) reminder.Reminder {
	return &SimpleReminder{
		message: msg,
	}
}
