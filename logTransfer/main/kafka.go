package main

import (
	"strings"
	"sync"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type KafkaClient struct {
	client sarama.Consumer
	addr   string
	topic  string
	wg     sync.WaitGroup
}

var (
	kafkaClient *KafkaClient
)

func initKafka(addr string, topic string) (err error) {
	logs.Debug("Start to init kafka.")
	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("Init kafka failed, Error :%v", err)
		return
	}
	kafkaClient = &KafkaClient{}
	kafkaClient.client = consumer
	kafkaClient.addr = addr
	kafkaClient.topic = topic
	logs.Debug("Successful initialized kafka.")
	return
}
