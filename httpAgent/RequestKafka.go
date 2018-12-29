package httpAgent

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"runtime/debug"
	"strings"
	"time"
)

const (
	_BROKER_LIST_ = `localhost:9092`
)
const (
	_LABEL_ = "[_httpAgent_]"
)

var (
	IS_DEBUG = false
	_PAUSE_  = false
)

func SetDebug(debug bool) {
	IS_DEBUG = debug
}

type KafkaAgent struct {
	flag           bool
	BrokerList     string
	TopicList      string
	SendTimeOut    time.Duration
	ReceiveTimeOut time.Duration
	AsyncProducer  sarama.AsyncProducer
}

func (this *KafkaAgent) SetKafkaConf(BrokerList, TopicList string, SendTimeOut time.Duration, ReceiveTimeOut time.Duration) bool {
	//只允许初始化一次
	if this.flag {
		return false
	}

	this.flag = true
	this.BrokerList = BrokerList
	this.TopicList = TopicList
	this.SendTimeOut = SendTimeOut
	this.ReceiveTimeOut = ReceiveTimeOut
	this.AsyncProducer = getProducer(this.BrokerList, this.SendTimeOut, true)
	if nil == this.AsyncProducer {

		return false
	}

	return this.Check()
}
func (this *KafkaAgent) Check() bool {
	if "" == this.BrokerList || "" == this.TopicList {
		return false
	}
	if 0 == this.SendTimeOut && 0 == this.ReceiveTimeOut {
		return false
	}
	return true
}
func (this *KafkaAgent) SendMessege(msg string) bool {
	defer func() {
		if e, ok := recover().(error); ok {
			log.Println(_LABEL_, "WARN: panic in %v", e)
			log.Println(_LABEL_, string(debug.Stack()))
			this.AsyncProducer.Close()
			this.AsyncProducer = getProducer(this.BrokerList, this.SendTimeOut, true)
		}
	}()
	if !this.Check() {

		return false
	}

	return asyncProducer(
		this.AsyncProducer,
		this.TopicList,
		msg,
	)
}

//=========================================================================
// asyncProducer 异步生产者
func AsyncProducer(kafka_list, topics, s string, timeout time.Duration) bool {
	if "" == kafka_list || "" == topics {
		return false
	}
	producer := getProducer(kafka_list, timeout, false)
	if nil == producer {
		return false
	}
	defer producer.Close()
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					if IS_DEBUG {
						log.Println(_LABEL_, err)
					}
					return
				} else {
					return
				}
			case <-success:
				return
			}
		}
	}(producer)

	return asyncProducer(producer, topics, s)
}
func asyncProducer(p sarama.AsyncProducer, topics, s string) bool {
	if nil == p {
		return false
	}

	msg := &sarama.ProducerMessage{
		Topic: topics,
		Value: sarama.ByteEncoder(s),
	}
	p.Input() <- msg

	if IS_DEBUG {
		fmt.Println(_LABEL_, msg)
	}
	return true
}
func getProducer(kafka_list string, timeout time.Duration, monitor bool) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = timeout
	producer, err := sarama.NewAsyncProducer(strings.Split(kafka_list, ","), config)
	if err != nil {
		if IS_DEBUG {
			log.Println(_LABEL_, err)
		}
	}
	if monitor {
		//消费状态消息,防止死锁
		go func(producer sarama.AsyncProducer) {
			if nil == producer {
				log.Println(_LABEL_, "getProducer() producer error!")

				return
			}
			errors := producer.Errors()
			success := producer.Successes()
			for {
				select {
				case err := <-errors:
					if err != nil {
						if IS_DEBUG {
							log.Println(_LABEL_, err)
						}
						continue
					} else {
						continue
					}
				case <-success:
					continue
				}
			}
		}(producer)
	}
	return producer
}
