package reminder

import (
	"time"

	"github.com/google/uuid"
)

const (
	createdStatusString = "CREATED"
	queuedStatusString  = "QUEUED"
	firedStatusString   = "FIRED"
	missedStatusString  = "MISSED"
)

var (
	msgStatusToString = map[MsgStatus]string{
		Created: createdStatusString,
		Queued:  queuedStatusString,
		Fired:   firedStatusString,
		Missed:  missedStatusString,
	}
)

// MsgStatus is the status that a message is at
type MsgStatus int

func (s MsgStatus) String() string {
	return msgStatusToString[s]
}

const (
	_ MsgStatus = iota
	// Created represents that a message has been created but not acted on yet since
	Created

	// Queued means the system has taken the message and placed it into memory; it _should_ be there
	Queued

	// Fired means that this message - to the extent that the system can promise - has been delivered
	Fired

	// Missed is applied to a message if for some reason we did not catch the message when we were supposed to
	// TODO: use this later
	Missed
)

// Reminder represents the basic reminder type that is used at the API layer
type Reminder struct {
	ID      string
	Message string
	To      string
	Created int64
	Moment  int64
	Status  MsgStatus
}

// New initializes a reminder; applies the ID and the Status. Clients SHOULD NOT directly use the reminder struct
func New(created, moment int64, message, to string) *Reminder {
	return &Reminder{
		ID:      uuid.New().String(),
		Created: created,
		Moment:  moment,
		Message: message,
		To:      to,
		Status:  Created,
	}
}

func (r *Reminder) IsAfter(t time.Duration) bool {
	var ttl = time.Duration(r.Moment - time.Now().Unix())

	if ttl > 0 && ttl < t {
		return true
	}

	return false
}

func (r *Reminder) CalcTTL() time.Duration {
	return time.Duration(r.Moment-time.Now().Unix()) * time.Second
}
