package rest

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/kelindar/binary"
	"github.com/scottshotgg/reminder/pkg/reminder"
	v1 "github.com/scottshotgg/reminder/pkg/reminder/v1"
	"github.com/scottshotgg/reminder/pkg/sender"
)

type Server struct {
	s  sender.Sender
	r  *redis.Client
	ch chan string
}

func Start(s sender.Sender) {
	var app = fiber.New()

	var srv = &Server{
		s:  s,
		ch: make(chan string, 1000),
	}

	// TODO: options
	srv.r = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go srv.crawl()
	var amount = 10

	// Start the workers
	for i := 0; i < amount; i++ {
		go srv.worker()
	}

	app.Get("/reminders", srv.getReminder)
	app.Post("/reminders", srv.createReminder)

	app.Listen(8080)
}

/*
 - Get all the keys
 - Loop through the keys
 - Set time.AfterFunc for each one under 5 minutes
*/
func (s *Server) process(key string) {
	fmt.Println("Got key:", key)

	// Grab the TTL of the given key
	ttl, err := s.r.TTL(key).Result()
	if err != nil {
		log.Fatalln("err:", err)
	}

	// If TTL is over 5 minutes, ignore it for now
	if ttl > 5*time.Minute {
		fmt.Println("Skipping ...")
		return
	}

	// Get the key
	contents, err := s.r.Get(key).Bytes()
	if err != nil {
		log.Fatalln("err:", err)
	}

	// Unmarshal the contents from the binary payload
	var r v1.V1
	err = binary.Unmarshal(contents, &r)
	if err != nil {
		log.Fatalln("err:", err)
	}

	fmt.Println("reminder:", r)

	// If its already been queued then skip it
	// TODO: put this in another Redis hash set or w/e
	if r.Queued {
		fmt.Println("Already queued")
		return
	}

	fmt.Println("Setting reminder ...")
	// Set a timer to fire the send
	s.after(ttl, &r, key)
}

func (s *Server) after(ttl time.Duration, r reminder.Reminder, key string) {
	time.AfterFunc(ttl, func() {
		var err = s.s.Send(r)
		if err != nil {
			// TODO: need to do something here
			log.Fatalln("err:", err)
		}

		// Delete the key after we are done with it
		// TODO: see if Redis will do this automatically when we are done with it
		_, err = s.r.Del(key).Result()
		if err != nil {
			log.Fatalln("err:", err)
		}
	})
}

func (s *Server) worker() {
	for key := range s.ch {
		s.process(key)
	}
}

func (s *Server) crawl() {
	for range time.NewTicker(1 * time.Minute).C {
		var keys, err = s.r.Keys("*").Result()
		if err != nil {
			log.Fatalln("err:", err)
		}

		for _, key := range keys {
			s.ch <- key
		}
	}
}

func (s *Server) getReminder(c *fiber.Ctx) {
	c.Send("Hello, World 👋!")
}

type (
	V1 struct {
		id      uuid.UUID
		Created int64
		After   int64
		Msg     string
		To      string
	}
)

func (s *Server) createReminder(c *fiber.Ctx) {
	// var randomNum = rand.Intn(6)
	// var randy = "reminder" + strconv.Itoa(randomNum)

	// var blob, err = json.Marshal(simp.New(randy))
	// if err != nil {
	// 	log.Fatalln("err:", err)
	// }

	// var cmd = s.r.Set(randy, blob, time.Duration(randomNum)*time.Minute)
	// if cmd.Err() != nil {
	// 	log.Fatalln("err:", cmd.Err().Error())
	// }

	var rem V1
	var err = c.BodyParser(&rem)
	if err != nil {
		log.Fatalln("err:", err)
	}

	var now = time.Now()

	if rem.Created == 0 {
		rem.Created = now.Unix()
	}

	rem.id, err = uuid.NewUUID()
	if err != nil {
		log.Fatalln("err:", err)
	}

	fmt.Println("rem:", rem)

	var v1Rem = v1.V1{
		Created: rem.Created,
		After:   time.Duration(rem.After) * time.Second,
		Msg:     rem.Msg,
		To:      rem.To,
	}

	// If its going to fire in less than 5 minutes then instantly queue it up
	if v1Rem.After < 5*time.Minute {
		v1Rem.Queued = true

		s.after(v1Rem.After, &v1Rem, rem.id.String())
	}
	// If not then insert it into Redis
	blob, err := binary.Marshal(v1Rem)
	if err != nil {
		log.Fatalln("err:", err)
	}

	fmt.Println(s.r.Set(rem.id.String(), blob, v1Rem.After).Result())
}
