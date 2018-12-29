package httpAgent

import (
	"fmt"
	"github.com/golang/glog"
	"os"
	"strconv"
	"sync"
	"time"
)

func GetDailyGetReqData(httpConf ConfSlice, kafkaconf KafkaAgent) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(httpConf.GetPath); i++ {
			if httpConf.GetPath[i].UpDatetime != "0" {
				continue
			}
			RequestGetData(httpConf.GetPath[i].Url, httpConf.GetPath[i].Command, kafkaconf)
		}
		now := time.Now()
		next := now.Add(time.Hour * 24)
		startHour, _ := strconv.Atoi(httpConf.Default.StartTime)
		next = time.Date(next.Year(), next.Month(), next.Day(), startHour, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}()
	wg.Wait()
}

func GetTimerGetReqData(httpConf ConfSlice, kafkaConf KafkaAgent) {
	var wg sync.WaitGroup
	for j := 0; j < len(httpConf.GetPath); j++ {
		if httpConf.GetPath[j].UpDatetime == "0" {
			continue
		}
		wg.Add(1)
		fmt.Println(httpConf.GetPath[j].UpDatetime)
		mm, _ := time.ParseDuration(httpConf.GetPath[j].UpDatetime)
		timer2 := time.NewTimer(mm)
		go func(index int) {
			defer wg.Done()
			RequestGetData(httpConf.GetPath[index].Url, httpConf.GetPath[index].Command, kafkaConf)
			<-timer2.C
		}(j)
	}
	wg.Wait()
}

func GetDailyPostReqData(httpConf ConfSlice, kafkaConf KafkaAgent) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < len(httpConf.PostPath); i++ {
			if httpConf.PostPath[i].UpDatetime != "0" {
				continue
			}
			RequestPostTimeData(httpConf.PostPath[i].Url, httpConf.PostPath[i].Command, httpConf.PostPath[i].StartPos, httpConf.PostPath[i].EndPos, httpConf.PostPath[i].UpDatetime, true, kafkaConf)
		}
		now := time.Now()
		next := now.Add(time.Hour * 24)
		startHour, _ := strconv.Atoi(httpConf.Default.StartTime)
		next = time.Date(next.Year(), next.Month(), next.Day(), startHour, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}()
	wg.Wait()
}

func GetTimerPostReqData(httpConf ConfSlice, kafkaConf KafkaAgent) {
	var wg sync.WaitGroup
	for k := 0; k < len(httpConf.PostPath); k++ {
		if httpConf.PostPath[k].UpDatetime == "0" {
			continue
		}
		wg.Add(1)
		// 通过startdate字段  来区分使用导入以前天数据
		if httpConf.PostPath[k].StartDate == "0" {
			mm, _ := time.ParseDuration(httpConf.PostPath[k].UpDatetime)

			//var updatetime time.Duration   =  time.Duration(uptime) * time.Second
			timer_post := time.NewTimer(mm)
			go func(index int) {
				defer wg.Done()
				fmt.Println(index)
				RequestPostTimeData(httpConf.PostPath[index].Url, httpConf.PostPath[index].Command, httpConf.PostPath[index].StartPos, httpConf.PostPath[index].EndPos, httpConf.PostPath[index].UpDatetime, false, kafkaConf)
				<-timer_post.C
			}(k)
		} else {
			mm, _ := time.ParseDuration(httpConf.PostPath[k].UpDatetime)
			timer_post := time.NewTimer(mm)
			go func(index int) {
				defer wg.Done()
				fmt.Println(index)
				RequestPostBeforeTimeData(httpConf.PostPath[index].Url, httpConf.PostPath[index].Command, httpConf.PostPath[index].StartPos, httpConf.PostPath[index].EndPos, httpConf.PostPath[index].UpDatetime, httpConf.PostPath[index].StartDate, kafkaConf)
				<-timer_post.C
				os.Exit(0)
				glog.Info("RequestPostBeforeTimeData  successful")
			}(k)
		}
	}
	wg.Wait()
}
