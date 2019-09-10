package cNatsHelper

import (
	"encoding/json"
	"errors"
	"github.com/nats-io/go-nats"
	"os"
	"time"
)

var (
	cnh *cNatsHelper
)

type SubHandler func(msg CMsg)

func InitAndLoadHelper(modName , addr string) error {
	if modName == "" || addr == "" { return errors.New("mod Name or address can't be empty") }

	if cnh == nil {
		conn , err := nats.Connect(addr)
		if err != nil { return err }

		cnh = &cNatsHelper{
			mod: modName,
			pid: int64(os.Getegid()),
			conn:  conn,
		}
	}
	return nil
}

type cNatsHelper struct {
	mod 	string
	pid 	int64
	conn  	*nats.Conn
}

func (c *cNatsHelper) NewMsg(mod , serve string , data []byte) CMsg {
	cMsg := CMsg{
		FInfo:     SerInfo{
			Mod: c.mod,
			Ser: "",
			PID: c.pid,
		},
		TInfo:     SerInfo{
			Mod: mod,
			Ser: serve,
			PID: 0,
		},
		TimeStamp: time.Now().Unix(),
		Data:      data,
	}
	return cMsg
}

// publish task
func (c *cNatsHelper) publish(mod , serve string , data []byte) error {
	cMsg := c.NewMsg(mod , serve , data)
	return c.conn.Publish(cMsg.GetToAddr() , cMsg.ToJSONBytes())
}

// do task
func (c *cNatsHelper) subscribe(serve string , handler nats.MsgHandler) error {
	subject := c.mod + "." + serve
	_ , err := c.conn.Subscribe(subject , handler)
	return err
}

// queue do task
func (c *cNatsHelper) queueSubscribe(serve , queue string , handler nats.MsgHandler) error {
	subject := c.mod + "." + serve
	_ , err := c.conn.QueueSubscribe(subject , queue , handler)
	return err
}

// sync request - response
func (c *cNatsHelper) requestSync(mod, serve string, data []byte, timeout time.Duration) (CMsg, error) {
	reqMsg := c.NewMsg(mod , serve , data)
	msg , err := c.conn.Request(reqMsg.GetToAddr() , reqMsg.ToJSONBytes() , timeout)
	if err != nil {
		return CMsg{} , err
	}

	var rspMsg CMsg
	err = json.Unmarshal(msg.Data , &rspMsg)
	return rspMsg , err
}



func Publish(mod , serve string , data []byte) error {
	if mod == "" || serve == "" {
		return errors.New("mod , serve or handler can not be nil")
	}

	return cnh.publish(mod , serve , data)
}

func Subscribe(serve string , handler SubHandler) error {
	if serve == "" || handler == nil {
		return errors.New("mod , serve or handler can not be nil")
	}

	h := func(msg *nats.Msg) {
		cMsg := CMsg{}
		_ = json.Unmarshal(msg.Data , &cMsg)
		handler(cMsg)
	}

	return cnh.subscribe(serve , h)
}

func QueueSubscribe(serve , queue string , handler SubHandler) error {
	if serve == "" || handler == nil {
		return errors.New("mod , serve , queue or handler can not be nil")
	}

	h := func(msg *nats.Msg) {
		cMsg := CMsg{}
		_ = json.Unmarshal(msg.Data , &cMsg)
		handler(cMsg)
	}

	return cnh.queueSubscribe(serve , queue , h)
}

func Request(mod , serve string , data []byte , timeout time.Duration) (CMsg , error) {
	if mod == "" || serve == "" || data == nil {
		return CMsg{} , errors.New("mod , serve , queue or handler can not be nil")
	}

	return cnh.requestSync(mod , serve , data , timeout)
}