package storage

import (
	"time"

	"github.com/scottshotgg/reminder/pkg/types"
)

type Storage interface {
	ListReminders() ([]string, error)

	GetReminder(key string) (*types.DBReminder, error)
	GetTTL(key string) (time.Duration, error)

	CreateReminder(r *types.DBReminder) error
	DeleteKey(key string) error
}
