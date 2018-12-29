package httpAgent

import (
	"encoding/json"
	"io/ioutil"
)
type DefalutOptionConf struct {
	StartTime   string	`json:"starttime"`
	KafkaBrokes string	`json:"kafkabrokes"`
	KafkaTopic  string	`json:"kafkatopic"`
}

type GetPathOptionConf struct {
	Url          string  `json:"url"`
	Command      string  `json:"command"`
	UpDatetime   string  `json:"updatetime"`
}

type PostPathOptionConf struct {
	Url          string  `json:"url"`
	StartDate	 string  `json:"startdate"`
	StartPos     string  `json:"startPos"`
	EndPos       string   `json:"endPos"`
	Command      string  `json:"command"`
	UpDatetime   string  `json:"updatetime"`
}

type ConfSlice struct {
	Default    DefalutOptionConf
	GetPath    []GetPathOptionConf
	PostPath   []PostPathOptionConf
}

func ReadHttpAgentConf(filename string) ([]byte,error){
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return data,err
	}
	return data,nil
}

func ParseHttpAgentConf(filedata []byte) (ConfSlice,error){//
	var s ConfSlice
	err:=json.Unmarshal([]byte(filedata), &s)
	return s,err
}
