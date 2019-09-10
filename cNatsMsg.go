package cNatsHelper

import "encoding/json"

type SerInfo struct {
	Mod string
	Ser string
	PID int64
}

type CMsg struct {
	FInfo 		SerInfo
	TInfo 		SerInfo
	TimeStamp 	int64
	Data		[]byte
}

func (cMsg *CMsg) ToJSONBytes() []byte {
	jsonBytes , err := json.Marshal(cMsg)
	if err != nil {
		return []byte("")
	}
	return jsonBytes
}

func (cMsg *CMsg) GetToAddr() string {
	return cMsg.TInfo.Mod + "." + cMsg.TInfo.Ser
}