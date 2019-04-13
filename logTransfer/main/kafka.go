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
	kafkaClient *KafkaClient = &KafkaClient{}
)

func initKafka(addr string, topic string) (err error) {
	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("Init kafka failed, Error :%v", err)
		return
	}

	kafkaClient.client = consumer
	kafkaClient.addr = addr
	kafkaClient.topic = topic
	return
}
