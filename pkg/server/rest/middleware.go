package rest

import (
	"errors"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
)

func (s *Server) parseReminderID(c *fiber.Ctx) {
	var (
		id  = c.Params(s.params.reminderID)
		err error
	)

	if len(id) != 36 {
		err = errors.New("not 36 characters")
	}

	if err == nil {
		_, err = uuid.Parse(id)
	}

	if err != nil {
		err = c.JSON(NewError("invalid reminder ID", err.Error()))
		if err != nil {
			c.Send(`{"error": %s}`, err.Error())
		}

		return
	}

	c.Next()
}
