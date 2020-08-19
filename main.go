package main

import (
	"errors"
	"flag"
	"log"

	"github.com/scottshotgg/reminder/pkg/sender"
	"github.com/scottshotgg/reminder/pkg/sender/printer"
	"github.com/scottshotgg/reminder/pkg/sender/sms/twilio"
	"github.com/scottshotgg/reminder/pkg/server/rest"

	// redis_storage "github.com/scottshotgg/reminder/pkg/storage/redis"
	mem_storage "github.com/scottshotgg/reminder/pkg/storage/mem"
)

const (
	workerAmount = 10
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	// TODO: make a cmd line flag for this
	// var s, err = printer.New()
	// if err != nil {
	// 	log.Fatalln("err:", err)
	// }

	var (
		sender sender.Sender
		err    error

		// Utilities
		smsProviderFlag = flag.String("sms-provider", "printer", "Use this to set twilio as the default text provider")

		// TODO: fill in usage
		// Twilio
		fromFlag  = flag.String("from", "", "")
		sidFlag   = flag.String("sid", "", "")
		tokenFlag = flag.String("token", "", "")

		// Redis
		redisURLFlag = flag.String("redis-url", "localhost:6379", "")
	)

	// TODO: make a flag for storage
	// Storage
	// var redisURLFlag = flag.String("redis-url", "localhost:6379", "")
	// _, err = net.ParseIP(*redisURLFlag)
	// if err != nil {
	// 	log.Fatalln("You must provide the `redis-url` flag")
	// }

	flag.Parse()

	switch *smsProviderFlag {
	case "twilio":
		if *fromFlag == "" {
			log.Fatalln("You must provide the `from` flag")
		}

		if *sidFlag == "" {
			log.Fatalln("You must provide the `sid` flag")
		}

		if *tokenFlag == "" {
			log.Fatalln("You must provide the `token` flag")
		}

		sender, err = twilio.New(*fromFlag, *sidFlag, *tokenFlag)

	case "printer":
		sender, err = printer.New()

	case "":
		err = errors.New("no `smsProvider` specified")

	default:
		err = errors.New("invalid `smsProvider` specified: " + *smsProviderFlag)
	}

	if err != nil {
		log.Fatalln("err:", err)
	}

	// store, err := redis_storage.New(*redisURLFlag)
	// if err != nil {
	// 	log.Fatalln("err:", err)
	// }

	_ = *redisURLFlag

	var store = mem_storage.New(nil)

	err = rest.Start(workerAmount, store, sender)
	if err != nil {
		log.Fatalln("err:", err)
	}
}
