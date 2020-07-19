package mem

import (
	"errors"
	"time"

	"github.com/scottshotgg/reminder/pkg/storage"
	"github.com/scottshotgg/reminder/pkg/types"
)

type Mem struct {
	store map[string]*types.DBReminder
}

func New(kv map[string]*types.DBReminder) storage.Storage {
	if kv == nil {
		kv = map[string]*types.DBReminder{}
	}

	return &Mem{
		store: kv,
	}
}

func (s *Mem) ListReminders() ([]string, error) {
	return []string{}, nil
}

func (s *Mem) GetReminder(key string) (*types.DBReminder, error) {
	var r, ok = s.store[key]
	if !ok {
		return nil, errors.New("not found")
	}

	return r, nil
}

func (s *Mem) GetTTL(key string) (time.Duration, error) {
	var r, err = s.GetReminder(key)
	if err != nil {
		return 0, err
	}

	var finalTime = time.Unix(r.Created, 0).Add(r.Until)

	return time.Now().Sub(finalTime), nil
}

func (s *Mem) CreateReminder(r *types.DBReminder) error {
	s.store[r.ID] = r

	return nil
}

func (s *Mem) DeleteKey(key string) error {
	delete(s.store, key)

	return nil
}
