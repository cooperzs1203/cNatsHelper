package cNatsHelper

import (
	"errors"
	"os"
)

var (
	serveName = ""
	servePID = -1
)

func InitServe(Name string) error {
	if Name == "" {
		return errors.New("ServeName can not be empty")
	}

	serveName = Name
	servePID = os.Getegid()

	// init nats.EncodedConn
	err := initNatsEC()
	if err != nil {
		return err
	}

	return nil
}