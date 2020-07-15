package main

import (
	"log"

	"github.com/scottshotgg/reminder/pkg/sender/printer"
	"github.com/scottshotgg/reminder/pkg/server/rest"
)

func main() {
	//	fmt.Println("sup world")
	log.SetFlags(log.Lmicroseconds)

	var s, err = printer.New()
	if err != nil {
		log.Fatalln("err:", err)
	}

	// send(s)
	rest.Start(s)

	// time.Sleep(1 * time.Minute)
}

// func send(s sender.Sender) {
// 	var stop = time.After(1 * time.Second)

// 	for range time.NewTicker(100 * time.Millisecond).C {
// 		select {
// 		case <-stop:
// 			return

// 		default:
// 			log.Println("Sending")
// 			s.Send(simp.New("yo fker do ur things"))
// 		}
// 	}
// }
