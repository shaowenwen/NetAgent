package httpAgent

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func RequestGetData(reqUrl string, msgtype string,kafkaConf KafkaAgent) {
	resp, err := http.Post(reqUrl,"application/x-www-form-urlencoded",nil)
	if err != nil {
		glog.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
	}
	cn_json, _ := simplejson.NewJson([]byte(body))
	cn_body, _ := cn_json.Array()
	for _, jsonVal := range cn_body {
		jsonArray := PkgKafkaJsonData(jsonVal, msgtype)
		kafkaConf.SendMessege(jsonArray)
	}
}

func RequestPostData(reqUrl string, msgtype string, startpos string, endpos string) {
	var i int
	i, _ = strconv.Atoi(startpos)
	for {
		resp, err := http.PostForm(reqUrl, url.Values{"startPos": {strconv.Itoa(i)}, "endPos": {endpos}}) //
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		fmt.Println(resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
			fmt.Println(err)
			continue
		}
		fmt.Println(string(body))
		cn_json, err := simplejson.NewJson([]byte(body))
		if err != nil {
			// handle error
			fmt.Println(err)
			fmt.Println("request end")
			break
		}
		cn_body, _ := cn_json.Array()
		//fmt.Println(cn_body)
		for _, jsonVal := range cn_body {
			jsonArray := PkgKafkaJsonData(jsonVal, msgtype)
			fmt.Println(jsonArray)
		}
		limit, _ := strconv.Atoi(endpos)
		i = i + limit
	}
}

func RequestPostTimeData(reqUrl string, msgtype string, startpos string, endpos string, updatetime string, dailyFlag bool,kafkaConf KafkaAgent) {
	var i int
	i, _ = strconv.Atoi(startpos)
	for {
		now := time.Now()
		mm, _ := time.ParseDuration(updatetime)
		var _endTime string
		var _startTime string
		if dailyFlag {
			before := now.Add(-time.Hour * 24)
			before = time.Date(before.Year(), before.Month(), before.Day(), 0, 0, 0, 0, before.Location())
			_startTime = before.Format("2006-01-02 03:04:05")
			_endTime = now.Format("2006-01-02") + " " + "00:00:00"

		} else {
			_endTime = now.Add(mm).Format("2006-01-02 03:04:05")
			_startTime = now.Format("2006-01-02 03:04:05")
		}
		glog.Info(_startTime,_endTime)

		resp, err := http.PostForm(reqUrl, url.Values{"startPos": {strconv.Itoa(i)}, "endPos": {endpos}, "startTime": {_startTime}, "endTime": {_endTime}}) //
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200{
			glog.Error(resp.StatusCode,resp.Status)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
			glog.Error(err)
			continue
		}

		cn_json, err := simplejson.NewJson([]byte(body))
		if err != nil {
			// handle error
			glog.Error(err)
			glog.Error("request end")
			break
		}
		cn_body, _ := cn_json.Array()
		if len(cn_body) == 0 {
			glog.Error("cn_body array null")
			break
		}

		for _, jsonVal := range cn_body {
			jsonArray := PkgKafkaJsonData(jsonVal, msgtype)
			kafkaConf.SendMessege(jsonArray)
		}
		limit, _ := strconv.Atoi(endpos)
		i = i + limit
	}
}

func RequestPostBeforeTimeData(reqUrl string, msgtype string, startpos string, endpos string, updatetime string, startdate string,kafkaConf KafkaAgent) {
	var i int
	i, _ = strconv.Atoi(startpos)
	for {
		now := time.Now()
		_startDate := startdate + " " + "00:00:00"
		resp, err := http.PostForm(reqUrl, url.Values{"startPos": {strconv.Itoa(i)}, "endPos": {endpos}, "startTime": {_startDate}, "endTime": {now.Format("2006-01-02 03:04:05")}}) //
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		fmt.Println(resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
			glog.Error(err)
			continue
		}
		fmt.Println(string(body))
		cn_json, err := simplejson.NewJson([]byte(body))
		if err != nil {
			// handle error
			glog.Error(err)
			glog.Error("request end")
			break
		}
		cn_body, _ := cn_json.Array()
		for _, jsonVal := range cn_body {
			jsonArray := PkgKafkaJsonData(jsonVal, msgtype)
			kafkaConf.SendMessege(jsonArray)
		}
		limit, _ := strconv.Atoi(endpos)
		i = i + limit
	}
}


