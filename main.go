package main

import (
	"NJBDNet/httpAgent"
	"flag"
	"github.com/golang/glog"
	"time"
)

func InitGlog(){
	flag.Set("alsologtostderr", "true") // 日志写入文件的同时，输出到stderr
	flag.Set("log_dir", ".")        // 日志文件保存目录
	flag.Set("v", "3")                  // 配置V输出的等级。
	flag.Parse()
}

func main() {
	InitGlog()
	defer glog.Flush()

	glog.Info("httpAgent start !!!")

	confData, err := httpAgent.ReadHttpAgentConf("conf/httpconfig.json")
	if err != nil {
		glog.Error("httpAgent conf do not exist")
		return
	}
	httpConf ,err:= httpAgent.ParseHttpAgentConf(confData)
	if err != nil {
		glog.Error("httpAgent json conf error")
		return
	}
	kafkaConf  := httpAgent.KafkaAgent{}
	setKafkaFlag:=kafkaConf.SetKafkaConf(httpConf.Default.KafkaBrokes,httpConf.Default.KafkaTopic,5 * time.Second,5 * time.Second)
	if !setKafkaFlag{
		glog.Error("set Kafka info error")
	}

		go func() {
			for{
				httpAgent.GetTimerGetReqData(httpConf,kafkaConf)
				time.Sleep(time.Second)
			}
		}()
		go func() {
			for{
				httpAgent.GetDailyGetReqData(httpConf,kafkaConf)
				time.Sleep(time.Second)
			}
		}()
		go func() {
			for{
				httpAgent.GetTimerPostReqData(httpConf,kafkaConf)
				time.Sleep(time.Second)
			}
		}()
		go func() {
			for{
				httpAgent.GetDailyPostReqData(httpConf,kafkaConf)
				time.Sleep(time.Second)
			}
		}()

	for{
		time.Sleep(10*time.Minute)
	}



	glog.Info("httpAgent end !!!")

}
