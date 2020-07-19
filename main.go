package main

import (
	"errors"
	"flag"
	"log"

	"github.com/scottshotgg/reminder/pkg/sender"
	"github.com/scottshotgg/reminder/pkg/sender/printer"
	"github.com/scottshotgg/reminder/pkg/sender/sms/twilio"
	"github.com/scottshotgg/reminder/pkg/server/rest"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	// TODO: make a cmd line flag for this
	// var s, err = printer.New()
	// if err != nil {
	// 	log.Fatalln("err:", err)
	// }

	var (
		s   sender.Sender
		err error
	)

	// Utilities
	var smsProviderFlag = flag.String("sms-provider", "printer", "Use this to set twilio as the default text provider")

	// TODO: fill in usage
	// Twilio
	var fromFlag = flag.String("from", "", "")
	var sidFlag = flag.String("sid", "", "")
	var tokenFlag = flag.String("token", "", "")

	// Redis
	var redisURLFlag = flag.String("redis-url", "localhost:6379", "")
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

		s, err = twilio.New(*fromFlag, *sidFlag, *tokenFlag)

	case "printer":
		s, err = printer.New()

	case "":
		err = errors.New("no `smsProvider` specified")

	default:
		err = errors.New("invalid `smsProvider` specified: " + *smsProviderFlag)
	}

	if err != nil {
		log.Fatalln("err:", err)
	}

	err = rest.Start(*redisURLFlag, s)
	if err != nil {
		log.Fatalln("err:", err)
	}
}
