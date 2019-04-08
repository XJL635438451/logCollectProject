package main

import (
    "time"
    "github.com/Shopify/sarama"
     "github.com/hpcloud/tail"
     "go_dev/LogCollectProject/common"
    "github.com/astaxie/beego/logs"
)

var (
    client sarama.SyncProducer
)

func InitKafka(addr string) (err error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Partitioner = sarama.NewRandomPartitioner
    config.Producer.Return.Successes = true

    client, err = sarama.NewSyncProducer([]string{addr}, config)
    if err != nil {
        err = common.ErrMsg("producer close.", err)
        return
    }
    return 
}

func SendDataToKafka(tails *tail.Tail) {
    msg := &sarama.ProducerMessage{}
    msg.Topic = appConfig.Topic

    var chanMsg *tail.Line
    var ok bool
    for {
        logs.Debug("start to send message to kafka...")
        chanMsg, ok = <-tails.Lines
        if !ok {
            logs.Error("tail file close reopen, filename:%s\n", tails.Filename)
            time.Sleep(100 * time.Millisecond)
            continue
        }
        
        msg.Value = sarama.StringEncoder(chanMsg.Text)
        pid, offset, err := client.SendMessage(msg)
        if err != nil {
            logs.Error("send message failed. Error: ", err) 
            return 
        }

        logs.Debug("pid:%v offset:%v\n", pid, offset)
    }
}
