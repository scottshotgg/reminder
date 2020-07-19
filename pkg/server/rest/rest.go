package rest

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber"
	"github.com/scottshotgg/reminder/pkg/reminder"
	v1 "github.com/scottshotgg/reminder/pkg/reminder/v1"
	"github.com/scottshotgg/reminder/pkg/sender"
	"github.com/scottshotgg/reminder/pkg/storage"
	redis_storage "github.com/scottshotgg/reminder/pkg/storage/redis"
	"github.com/scottshotgg/reminder/pkg/types"
)

type Server struct {
	s       sender.Sender
	storage storage.Storage
	ch      chan string
}

func Start(redisURL string, s sender.Sender) error {
	var app = fiber.New()

	var srv = &Server{
		s:  s,
		ch: make(chan string, 1000),
	}

	var store, err = redis_storage.New(redisURL)
	if err != nil {
		return err
	}

	srv.storage = store

	go srv.crawl()
	var amount = 10

	// Start the workers
	for i := 0; i < amount; i++ {
		go srv.worker()
	}

	app.Get("/reminders", srv.getReminder)
	app.Post("/reminders", srv.createReminder)

	return app.Listen(8080)
}

/*
 - Get all the keys
 - Loop through the keys
 - Set time.AfterFunc for each one under 5 minutes
*/
func (s *Server) process(key string) error {
	// Grab the TTL of the given key
	ttl, err := s.storage.GetTTL(key)
	if err != nil {
		log.Fatalln("err:", err)
	}

	// If TTL is over 5 minutes, ignore it for now
	if ttl > 5*time.Minute {
		fmt.Println("Skipping ...")
		return nil
	}

	// Get the key
	r, err := s.storage.GetReminder(key)
	if err != nil {
		log.Fatalln("err:", err)
		return err
	}

	fmt.Println("reminder:", r)

	// If its already been queued then skip it
	// TODO: put this in another Redis hash set or w/e
	if r.Queued {
		fmt.Println("Already queued")
		return nil
	}

	fmt.Println("Setting reminder ...")

	// Set a timer to fire the send
	s.remind(ttl, v1.FromDB(r))

	return nil
}

// TODO: calculating the TTL needs to be more calculated, current time, etc
func (s *Server) remind(ttl time.Duration, r reminder.Reminder) {
	time.AfterFunc(ttl, func() {
		var err = s.s.Send(r)

		if err != nil {
			// TODO: need to do something here
			log.Fatalln("err:", err)
		}

		// TODO: we can probably optimize this by batching them
		// Delete the key after we are done with it
		err = s.storage.DeleteKey(r.GetID())
		if err != nil {
			log.Fatalln("err:", err)
		}
	})
}

func (s *Server) worker() {
	var err error

	for key := range s.ch {
		err = s.process(key)
		if err != nil {
			// TODO: do something here

			log.Fatalln("err:", err)
		}
	}
}

func (s *Server) crawl() {
	for range time.NewTicker(1 * time.Minute).C {
		var keys, err = s.storage.ListReminders()
		if err != nil {
			log.Fatalln("err:", err)

			// TODO: need to do something to decode the error and act
			continue
		}

		for _, key := range keys {
			s.ch <- key
		}
	}
}

func (s *Server) getReminder(c *fiber.Ctx) {
	c.Send("Hello, World 👋!")
}

func (s *Server) createReminder(c *fiber.Ctx) {
	var (
		rem V1
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

	var after = time.Duration(rem.After) * time.Second

	r, err := v1.New(rem.Created, after, rem.Message, rem.To)
	if err != nil {
		// TODO: return something
		log.Fatalln("err:", err)
	}

	// If its going to fire in less than 5 minutes then instantly queue it up
	if after < 5*time.Minute {
		r.SetQueued(true)

		// Start the reminder
		s.remind(after, r)
	}

	err = s.storage.CreateReminder(types.ToDB(r))
	if err != nil {
		// TODO: return something
		log.Fatalln("err:", err)
	}
}
