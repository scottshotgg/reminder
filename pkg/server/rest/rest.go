package rest

import (
	"github.com/gofiber/fiber"
	"github.com/scottshotgg/reminder/pkg/sender"
	"github.com/scottshotgg/reminder/pkg/storage"
)

const (
	reminderIDParam = "reminderID"
)

type params struct {
	reminderID string
}

type Server struct {
	sender  sender.Sender
	storage storage.Storage
	ch      chan string
	params  *params
}

func Start(workerAmount int, st storage.Storage, se sender.Sender) error {
	var (
		app = fiber.New()

		s = &Server{
			storage: st,
			sender:  se,
			ch:      make(chan string, 1000),
			params: &params{
				reminderID: reminderIDParam,
			},
		}
	)

	scrape(s, 10)

	// TODO: need to add account stuff in here
	// TODO: add paging later

	var (
		root = "/reminders"
		byID = root + "/:" + s.params.reminderID
	)

	// TODO: finish this and add query param for status to this one, list, and count
	// app.Get(root+"/from/:start/to/:end", s.getRemindersBetween)
	// app.Get(root+"/count", s.getReminderCount)

	// TODO: check less than zero on the creates and updates

	app.Get(root, s.listReminders)
	app.Get(byID, s.parseReminderID, s.getReminder)

	app.Post(root+"/at", s.createReminderAt)
	app.Post(root+"/after", s.createReminderAfter)

	app.Put(byID+"/at", s.parseReminderID, s.updateReminderAt)
	app.Put(byID+"/after", s.parseReminderID, s.updateReminderAfter)

	app.Delete(byID, s.parseReminderID, s.deleteReminder)

	return app.Listen(8080)
}
