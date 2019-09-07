package cNatsHelper

import (
	"fmt"
	"github.com/nats-io/go-nats"
	"log"
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

	log.Printf("Normal cMsg : %+v" , cMsg)

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

	log.Printf("Request cMsg : %+v" , cMsg)

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
		replyProfix := strings.Join([]string{reqCMsg.GetServeName() , reqCMsg.GetServePID() , _REQUEST_REPLY_} , ".")
		if strings.HasPrefix(reply, replyProfix) {
			reply = strings.Replace(reply , replyProfix + "." , "" , 1)
		}
	}

	cMsg.Msg.Subject = reply

	log.Printf("Reply cMsg : %+v" , cMsg)

	return cMsg
}

func buildSubject(subject string) string {
	if subject == "" {
		return ""
	}

	// use serveName.servePID to be the module part
	if serveName != "" && servePID != -1 {
		subject = strings.Join([]string{serveName , fmt.Sprintf("%d" , servePID) , subject} , ".")
	}

	log.Println("subject : " , subject)

	return subject
}

func buildQueue(queue string) string {
	if queue == "" {
		return ""
	}

	// use serveName_servePID to be the profix part
	if serveName != "" && servePID != -1 {
		queue = strings.Join([]string{serveName , fmt.Sprintf("%d" , servePID) , queue} , "_")
	}

	log.Println("queue : " , queue)

	return queue
}

func buildReply() string {
	var reply string

	// use serveName.servePID to be the profix part
	if serveName != "" && servePID != -1 {
		reply = strings.Join([]string{serveName , fmt.Sprintf("%d" , servePID) , _REQUEST_REPLY_} , ".")
	}

	log.Println("reply : " , reply)

	return reply
}