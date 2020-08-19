package rest

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"github.com/scottshotgg/reminder/pkg/reminder"
	"github.com/scottshotgg/reminder/pkg/types"
)

// TODO: rework this and make these handlers separate from the fiber function

func handleErrorRes(c *fiber.Ctx, err error) {

}

func handleSuccessRes(c *fiber.Ctx, value interface{}) {

}

func handleRes(c *fiber.Ctx, err error, value interface{}) {
	if err != nil {
		handleErrorRes(c, err)
	} else {
		handleSuccessRes(c, value)
	}
}

func (s *Server) listReminders(c *fiber.Ctx) {
	var keys, err = s.storage.ListReminders()
	if err != nil {
		err = c.JSON(NewError("", err.Error()))
		if err != nil {
			c.Send(`{"error": %s}`, err.Error())
		}

		return
	}

	err = c.JSON(ListV1Res{
		Reminders: keys,
	})

	if err != nil {
		log.Fatalln("err:", err)
	}
}

func (s *Server) getReminder(c *fiber.Ctx) {
	var id = c.Params("reminderID")

	r, err := s.storage.GetReminder(id)
	if err != nil {
		log.Fatalln("err:", err)
	}

	err = c.JSON(r)
	if err != nil {
		log.Fatalln("err:", err)
	}
}

func (s *Server) createReminderAt(c *fiber.Ctx) {
	var (
		rem CreateAtV1Req
		err = c.BodyParser(&rem)
	)

	if err != nil {
		// TODO: return something
		log.Fatalln("err:", err)
	}

	if rem.Created == 0 {
		rem.Created = time.Now().Unix()
	}

	fmt.Println("rem:", rem)

	var r = reminder.New(rem.Created, rem.Moment, rem.Message, rem.To)

	err = s.checkAndSet(c, r)
	if err != nil {
		// TODO: put something here later
		log.Fatalln("err:", err)
	}
}

func (s *Server) createReminderAfter(c *fiber.Ctx) {
	var (
		rem CreateAfterV1Req
		err = c.BodyParser(&rem)
	)

	if err != nil {
		// TODO: return something
		log.Fatalln("err:", err)
	}

	if rem.Created == 0 {
		rem.Created = time.Now().Unix()
	}

	fmt.Println("rem:", rem)

	var r = reminder.New(rem.Created, time.Now().Unix()+rem.After, rem.Message, rem.To)

	err = s.checkAndSet(c, r)
	if err != nil {
		// TODO: put something here later
		log.Fatalln("err:", err)
	}
}

func (s *Server) updateReminderAt(c *fiber.Ctx) {
	var id = c.Params("reminderID")

	rdb, err := s.storage.GetReminder(id)
	if err != nil {
		// TODO: put something here later
		log.Fatalln("err:", err)
	}

	// TODO: fill this out
	switch rdb.Status {

	}

	var (
		r   = dbToAPIReminder(rdb)
		rem CreateAtV1Req
	)

	err = c.BodyParser(&rem)
	if err != nil {
		// TODO: return something
		log.Fatalln("err:", err)
	}

	if rem.Created == 0 {
		r.Created = time.Now().Unix()
	}

	// TODO: consider making all of this in the middleware
	if rem.Moment < 0 {
		log.Fatalln("err: reminder about negative time would be weird ...")
	}

	if rem.Moment < rem.Created {
		log.Fatalln("err: i can't remind you about the past")
	}

	r.Moment = rem.Moment
	r.Message = rem.Message
	r.To = rem.To

	err = s.checkAndSet(c, r)
	if err != nil {
		// TODO: put something here later
		log.Fatalln("err:", err)
	}
}

func (s *Server) updateReminderAfter(c *fiber.Ctx) {
	var id = c.Params("reminderID")

	rdb, err := s.storage.GetReminder(id)
	if err != nil {
		// TODO: put something here later
		log.Fatalln("err:", err)
	}

	var r = dbToAPIReminder(rdb)

	r.ID = id

	var rem CreateAtV1Req

	err = c.BodyParser(&rem)
	if err != nil {
		// TODO: return something
		log.Fatalln("err:", err)
	}

	if rem.Created == 0 {
		r.Created = time.Now().Unix()
	}

	r.Moment = time.Now().Unix() + rem.Moment
	r.Message = rem.Message
	r.To = rem.To

	err = s.checkAndSet(c, r)
	if err != nil {
		// TODO: put something here later
		log.Fatalln("err:", err)
	}
}

func (s *Server) deleteReminder(c *fiber.Ctx) {
	var id = c.Params("reminderID")

	r, err := s.storage.GetReminder(id)
	if err != nil {
		log.Fatalln("err:", err)
	}

	err = s.storage.DeleteKey(id)
	if err != nil {
		log.Fatalln("err:", err)
	}

	err = c.JSON(r)
	if err != nil {
		log.Fatalln("err:", err)
	}
}

func (s *Server) checkAndSet(c *fiber.Ctx, r *reminder.Reminder) error {
	// If its going to fire in less than 5 minutes then instantly queue it up
	if r.IsAfter(5 * time.Minute) {
		r.Status = reminder.Queued

		// Start the reminder
		s.remind(r)
	}

	var err = s.storage.CreateReminder(types.ToDB(r))
	if err != nil {
		return err
	}

	err = c.JSON(r)
	if err != nil {
		return err
	}

	return nil
}
