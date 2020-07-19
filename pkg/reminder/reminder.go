package reminder

import (
	"time"
)

var (
	msgStatusToString = map[MsgStatus]string{
		Ready:  "READY",
		Queued: "QUEUED",
		Fired:  "FIRED",
		Missed: "MISSED",
	}
)

type MsgStatus int

func (s MsgStatus) String() string {
	return msgStatusToString[s]
}

const (
	_ MsgStatus = iota
	Ready
	Queued
	Fired

	// TODO: use this later
	Missed
)

type Reminder interface {
	GetID() string
	GetMessage() string
	GetTo() string
	GetCreated() int64
	GetMoment() int64
	GetStatus() MsgStatus

	IsReady() bool
	IsQueued() bool
	IsFired() bool

	SetQueued(q bool)

	CalcTTL() time.Duration
	IsAfter(t time.Duration) bool
}
