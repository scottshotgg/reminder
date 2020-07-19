package v1

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/scottshotgg/reminder/pkg/reminder"
	"github.com/scottshotgg/reminder/pkg/types"
)

type V1 struct {
	ID      string
	Created int64
	Until   time.Duration
	Message string
	Fired   bool
	Queued  bool
	To      string
}

func FromDB(r *types.DBReminder) *V1 {
	return &V1{
		ID:      r.ID,
		Created: r.Created,
		Until:   r.Until,
		Message: r.Message,
		Fired:   r.Fired,
		Queued:  r.Queued,
		To:      r.To,
	}
}

func New(created int64, until time.Duration, message, to string) (reminder.Reminder, error) {
	var id, err = uuid.NewUUID()
	if err != nil {
		log.Fatalln("err:", err)
		return nil, err
	}

	return &V1{
		ID:      id.String(),
		Created: created,
		Until:   until,
		Message: message,
		To:      to,
	}, nil
}

func (v *V1) GetID() string {
	return v.ID
}

func (v *V1) GetCreated() int64 {
	return v.Created
}

func (v *V1) GetUntil() time.Duration {
	return v.Until
}

func (v *V1) IsAfter(t time.Duration) bool {
	if v.Until < t {
		return true
	}

	return false
}

func (v *V1) GetMessage() string {
	return v.Message
}

func (v *V1) GetTo() string {
	return v.To
}

func (v *V1) GetQueued() bool {
	return v.Queued
}

func (v *V1) GetFired() bool {
	return v.Fired
}

func (v *V1) SetQueued(q bool) {
	v.Queued = q
}

// func (v V1) MarshalBinary() ([]byte, error) {
// 	return binary.Marshal(v)
// }

// func (v V1) UnmarshalBinary(data []byte) error {
// 	return binary.Unmarshal(data, v)
// }
