package cNatsHelper

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"strings"
)

const (
	_REQUEST_REPLY_ = "REQUEST_REPLY"
)

type ServeInfo struct {
	Name string
	PID int64
}

// CNatsMsg : For basic communication between services
// Msg : Communication underlying information
// serveInfo : public serve information
type CNatsMsg struct {
	ServeInfo ServeInfo
	Msg nats.Msg
}

func (cMsg *CNatsMsg) GetServeName() string {
	serveName := ""
	if cMsg == nil {
		return serveName
	}

	serveName = cMsg.ServeInfo.Name
	return serveName
}

func (cMsg *CNatsMsg) GetServePID() string {
	servePID := ""
	if cMsg == nil {
		return servePID
	}

	servePID = fmt.Sprintf("%d" , cMsg.ServeInfo.PID)
	return servePID
}

func (cMsg *CNatsMsg) GetSubject() string {
	subject := ""
	if cMsg == nil {
		return subject
	}

	subject = cMsg.Msg.Subject
	return subject
}

func (cMsg *CNatsMsg) GetMsgData() string {
	data := ""
	if cMsg == nil {
		return data
	}

	data = string(cMsg.Msg.Data)
	return data
}

func (cMsg *CNatsMsg) GetMsgDataBytes() []byte {
	var data []byte
	if cMsg == nil {
		return data
	}

	data = cMsg.Msg.Data
	return data
}

func (cMsg *CNatsMsg) GetReply() string {
	reply := ""
	if cMsg == nil {
		return reply
	}

	reply = cMsg.Msg.Reply
	return reply
}

func (cMsg *CNatsMsg) NeedReply() bool {
	if cMsg == nil {
		return false
	}

	return !(cMsg.Msg.Reply == "")
}

func (cMsg *CNatsMsg) GetReplyPrefix() string {
	replyPrefix := ""
	if cMsg == nil {
		return replyPrefix
	}

	replyPrefix = strings.Join([]string{cMsg.GetServeName() , cMsg.GetServePID() , _REQUEST_REPLY_} , ".")

	return replyPrefix
}




// make normal cNatsMsg for publish
func NormalCMsg(subject string, data []byte) CNatsMsg {
	if subject = buildSubject(subject); subject == "" {
		return CNatsMsg{}
	}

	cMsg := CNatsMsg{
		ServeInfo: ServeInfo{
			Name: serveName,
			PID:  int64(servePID),
		},
		Msg:       nats.Msg{
			Subject: subject,
			Reply: "",
			Data: data,
			Sub: nil,
		},
	}

	return cMsg
}

// a cNatsMsg for request
func RequestCMsg(subject string, data []byte) CNatsMsg {
	if subject = buildSubject(subject); subject == "" {
		return CNatsMsg{}
	}

	reply := buildReply()

	cMsg := CNatsMsg{
		ServeInfo: ServeInfo{
			Name: serveName,
			PID:  int64(servePID),
		},
		Msg:       nats.Msg{
			Subject: subject,
			Reply: reply,
			Data: data,
			Sub: nil,
		},
	}

	return cMsg
}

func ReplyCMsg(reqCMsg CNatsMsg , replyData []byte) CNatsMsg {
	cMsg := CNatsMsg{
		ServeInfo: ServeInfo{
			Name: serveName,
			PID:  int64(servePID),
		},
		Msg:       nats.Msg{
			Subject: "",
			Reply: "",
			Data: replyData,
			Sub: nil,
		},
	}

	if !reqCMsg.NeedReply() {
		return cMsg
	}

	reply := reqCMsg.GetReply()

	if reqCMsg.ServeInfo.Name != "" && reqCMsg.ServeInfo.PID != -1 {
		if strings.HasPrefix(reply, reqCMsg.GetReplyPrefix()) {
			reply = strings.Replace(reply , reqCMsg.GetReplyPrefix() , "" , 1)
		}
	}

	cMsg.Msg.Subject = reply

	return cMsg
}

