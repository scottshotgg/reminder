package printer

import (
	"log"

	"github.com/scottshotgg/reminder/pkg/reminder"
	"github.com/scottshotgg/reminder/pkg/sender"
)

// Printer implements Sender but is just a mock/test to just print
type Printer struct{}

func New() (sender.Sender, error) {
	return &Printer{}, nil
}

// Send is implemented very simply using time.AfterFunc
func (p *Printer) Send(r *reminder.Reminder) error {
	log.Println("MESSAGE:", r.Message)

	return nil
}
