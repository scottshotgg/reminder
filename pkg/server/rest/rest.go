package rest

import (
	"github.com/gofiber/fiber"
	"github.com/scottshotgg/reminder/pkg/sender"
)

func Start(s sender.Sender) {
	var app = fiber.New()

	app.Get("/reminders", getReminder)
	app.Post("/reminders", createReminder)

	app.Listen(8080)
}

func getReminder(c *fiber.Ctx) {
	c.Send("Hello, World 👋!")
}

func createReminder(c *fiber.Ctx) {
	c.Send("Hello, World 👋!")
}
