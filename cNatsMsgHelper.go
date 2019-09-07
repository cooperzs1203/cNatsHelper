package cNatsHelper

import (
	"fmt"
	"strings"
)

// subject should use serveName.servePID to be the prefix
func buildSubject(subject string) string {
	if subject == "" {
		return ""
	}

	// use serveName.servePID to be the module part
	if serveName != "" && servePID != -1 {
		subject = strings.Join([]string{serveName , fmt.Sprintf("%d" , servePID) , subject} , ".")
	}

	return subject
}

// queue should use serveName_servePID to be the prefix
func buildQueue(queue string) string {
	if queue == "" {
		return ""
	}

	// use serveName_servePID to be the profix part
	if serveName != "" && servePID != -1 {
		queue = strings.Join([]string{serveName , fmt.Sprintf("%d" , servePID) , queue} , "_")
	}

	return queue
}

// reply subject should use serveName.servePID to be the prefix
func buildReply() string {
	var reply string

	// use serveName.servePID to be the profix part
	if serveName != "" && servePID != -1 {
		reply = strings.Join([]string{serveName , fmt.Sprintf("%d" , servePID) , _REQUEST_REPLY_} , ".")
	}

	return reply
}