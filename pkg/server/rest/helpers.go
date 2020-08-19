package rest

import (
	"log"
	"time"

	"github.com/scottshotgg/reminder/pkg/reminder"
	"github.com/scottshotgg/reminder/pkg/types"
)

// TODO: calculating the TTL needs to be more calculated, current time, etc
func (s *Server) remind(r *reminder.Reminder) {
	time.AfterFunc(r.CalcTTL(), func() {
		var err = s.sender.Send(r)
		if err != nil {
			// TODO: need to do something here
			log.Fatalln("err:", err)
		}

		// TODO: we can probably optimize this by batching them
		// Delete the key after we are done with it
		err = s.storage.DeleteKey(r.ID)
		if err != nil {
			log.Fatalln("err:", err)
		}
	})
}

// dbToAPIReminder is a mapper that converts a persisted reminder into an API model
func dbToAPIReminder(r *types.DBReminder) *reminder.Reminder {
	return &reminder.Reminder{
		ID:      r.ID,
		Message: r.Message,
		To:      r.To,
		Moment:  r.Moment,
		Created: r.Created,
		Status:  r.Status,
	}
}
