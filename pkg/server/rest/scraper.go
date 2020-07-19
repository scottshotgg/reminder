package rest

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/scottshotgg/reminder/pkg/reminder"
	v1 "github.com/scottshotgg/reminder/pkg/reminder/v1"
)

// TODO: make this an interface in its own package so we can test different scraping algorithms

func scrape(s *Server, workerAmount int) {
	go s.crawl()

	// Start the workers
	for i := 0; i < workerAmount; i++ {
		go s.worker()
	}
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

	// If its already been queued then skip it
	// TODO: put this in another Redis hash set or w/e
	switch r.Status {
	case reminder.Ready:
		// Set a timer to fire the send
		s.remind(v1.FromDB(r))

	case
		reminder.Queued,
		reminder.Fired,
		reminder.Missed:

		fmt.Println("Already " + r.Status.String())

	default:
		return errors.New("invalid message status: " + r.Status.String())
	}

	return nil
}
