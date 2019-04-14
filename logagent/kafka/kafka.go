package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

var (
	client sarama.SyncProducer
)

func InitKafka(addr string) (err error) {
	logs.Debug("Start to initialize kafka.")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err = sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logs.Error("init kafka producer failed, err:", err)
		return
	}

	logs.Debug("Successfully initialized kafka.")
	return
}

func SendToKafka(data, topic string) (err error) {
	logs.Debug("Start to send message [msg: %s, topic: %s] to kafka.", data, topic)
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	pid, offset, err := client.SendMessage(msg)
	// _, _, err = client.SendMessage(msg)
	if err != nil {
		err = fmt.Errorf("Send message failed,  data:%v, topic:%v, Error: %v", data, topic, err)
		return
	}

	logs.Debug("Successful sent data to kafka, pid:%v offset:%v, topic:%v\n", pid, offset, topic)
	return
}
