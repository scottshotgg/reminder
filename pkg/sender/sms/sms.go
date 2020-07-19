package sms

import (
	"errors"
	"log"
	"regexp"
)

const (
	naPhoneNumRegEx = "(?:^|[^0-9])(1[34578][0-9]{9})(?:$|[^0-9])"
)

func init() {
	var _, err = regexp.Compile(naPhoneNumRegEx)
	if err != nil {
		log.Fatalln("err:", err)
	}
}

// type FromNumber struct {
// 	CountryCode string
// 	AreaCode    string
// 	Prefix      string
// 	Subscriber  string
// }

// // TODO: do some checking on the from number here and a parser from string
// func NewFrom(c, a, p, s string) (*FromNumber, error) {
// 	return &FromNumber{
// 		CountryCode: c,
// 		AreaCode:    a,
// 		Prefix:      p,
// 		Subscriber:  s,
// 	}, nil
// }

// func (f *FromNumber) String() string {
// 	return "+" + f.CountryCode + f.AreaCode + f.Prefix + f.Subscriber
// }

// TODO: only only works for North American numbers
func ValidateFrom(from string) error {
	var ok, err = regexp.MatchString(naPhoneNumRegEx, from)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Number did not pass regex")
	}

	return nil
}
