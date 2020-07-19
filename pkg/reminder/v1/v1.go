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
	Message string
	To      string
	Created int64
	Moment  int64
	Status  reminder.MsgStatus
}

func FromDB(r *types.DBReminder) *V1 {
	return &V1{
		ID:      r.ID,
		Message: r.Message,
		To:      r.To,
		Moment:  r.Moment,
		Created: r.Created,
		Status:  r.Status,
	}
}

func New(created, moment int64, message, to string) (reminder.Reminder, error) {
	var id, err = uuid.NewUUID()
	if err != nil {
		log.Fatalln("err:", err)
		return nil, err
	}

	return &V1{
		ID:      id.String(),
		Message: message,
		Created: created,
		Moment:  moment,
		To:      to,
	}, nil
}

func (v *V1) GetID() string {
	return v.ID
}

func (v *V1) GetMessage() string {
	return v.Message
}

func (v *V1) GetTo() string {
	return v.To
}

func (v *V1) GetCreated() int64 {
	return v.Created
}

func (v *V1) GetMoment() int64 {
	return v.Moment
}

func (v *V1) GetStatus() reminder.MsgStatus {
	return v.Status
}

func (v *V1) IsReady() bool {
	return v.Status == reminder.Ready
}

func (v *V1) IsQueued() bool {
	return v.Status == reminder.Queued
}

func (v *V1) IsFired() bool {
	return v.Status == reminder.Fired
}

func (v *V1) SetQueued(q bool) {
	v.Status = reminder.Queued
}

func (v *V1) CalcTTL() time.Duration {
	return time.Duration(v.Moment-time.Now().Unix()) * time.Second
}

func (v *V1) IsAfter(t time.Duration) bool {
	var ttl = time.Duration(v.Moment - time.Now().Unix())

	if ttl > 0 && ttl < t {
		return true
	}

	return false
}

// func (v V1) MarshalBinary() ([]byte, error) {
// 	return binary.Marshal(v)
// }

// func (v V1) UnmarshalBinary(data []byte) error {
// 	return binary.Unmarshal(data, v)
// }
