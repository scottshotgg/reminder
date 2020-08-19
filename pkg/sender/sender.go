package sender

import (
	"github.com/scottshotgg/reminder/pkg/reminder"
)

type Sender interface {
	Send(r *reminder.Reminder) error
}
