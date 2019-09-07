package cNatsHelper

import (
	"encoding/json"
	"errors"
	"github.com/nats-io/go-nats"
	"log"
	"time"
)

var (
	ec *nats.EncodedConn
)

const (
	_EMPTY_STRING_ = ""

	_NATS_REQUEST_TIMEOUT_ = time.Duration(15) * time.Second
)

func initNatsEC() error {
	addr := GlobalConfig().GetNatsServeAddr()
	conn , err := nats.Connect(addr)
	if err != nil {
		log.Printf("nats.Connect %s error : %s" , addr , err.Error())
		return err
	}

	eConn , err := nats.NewEncodedConn(conn , nats.JSON_ENCODER)
	if err != nil {
		log.Printf("nats.NewEdcodedConn %s error : %s" , addr , err.Error())
		return err
	}

	ec = eConn

	return nil
}

// publish message without reply
func Publish(cMsg CNatsMsg) error {
	return ec.Publish(cMsg.GetSubject() , cMsg)
}

// subscribe subject without queue
func Subscribe(subject string , handler subHandler) error {
	return subscribe(subject , _EMPTY_STRING_ , handler)
}

// subscribe subject with queue
func QueueSubscribe(subject, queue string, handler subHandler) error {
	return subscribe(subject , queue , handler)
}

// base subscribe function
func subscribe(subject, queue string, handler subHandler) error {
	if subject = buildSubject(subject); subject == "" || handler == nil {
		return errors.New("subject or handler can not be nil")
	}

	blockHandler := func(msg *nats.Msg) {
		cMsg := CNatsMsg{}
		err := json.Unmarshal(msg.Data , &cMsg)
		if err != nil {
			panic("nats.Msg Unmarshal to CNatsMsg error : " + err.Error())
		}

		handler(cMsg)
	}

	// if queue is not empty , that's mean it's a queue subscribe
	var err error
	if queue == "" {
		_ , err = ec.Subscribe(subject , blockHandler)
	} else {
		queue = buildQueue(queue)
		_ , err = ec.QueueSubscribe(subject , queue , blockHandler)
	}
	if err != nil {
		return err
	}

	return nil
}

// TODO:Use publish and subscribe to make sync request-reply and async request-reply
// request for sync reply
func Request(reqMsg CNatsMsg) (CNatsMsg , error) {
	var rspMsg nats.Msg
	err := ec.Request(reqMsg.GetSubject() , reqMsg , &rspMsg , _NATS_REQUEST_TIMEOUT_)
	if err != nil {
		return CNatsMsg{} , err
	}

	var rspCMsg CNatsMsg
	err = json.Unmarshal(rspMsg.Data , &rspCMsg)
	if err != nil {
		return CNatsMsg{} , err
	}

	return rspCMsg , nil
}

// reply for request
func Reply(rspMsg CNatsMsg) error {
	return Publish(rspMsg)
}

func ListenAndServe(subject string, handler subHandler) error  {
	if subject = buildSubject(subject); subject == "" || handler == nil {
		return errors.New("subject or handler can not be nil")
	}

	blockHandler := func(msg *nats.Msg) {
		cMsg := CNatsMsg{}
		err := json.Unmarshal(msg.Data , &cMsg)
		if err != nil {
			panic("nats.Msg Unmarshal to CNatsMsg error : " + err.Error())
		}

		cMsg.Msg.Reply = cMsg.Msg.Reply + "." + msg.Reply

		handler(cMsg)
	}

	var err error
	_ , err = ec.Subscribe(subject , blockHandler)
	if err != nil {
		return err
	}

	return nil
}






type subHandler func(msg CNatsMsg)
