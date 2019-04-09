package main

import (
    "time"
    "github.com/Shopify/sarama"
     "github.com/hpcloud/tail"
     "go_dev/LogCollectProject/common"
    "github.com/astaxie/beego/logs"
    "sync"
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

func GetDataFromKafka() {
    var wg sync.WaitGroup
    consumer, err := sarama.NewConsumer([]string{appConfig.Kafka_addr}, nil)
    if err != nil {
        logs.Error("consumer connect error:", err)
        return
    }
    logs.Debug("consumer connnect success...")
    defer consumer.Close()

    for {
        partitions, err := consumer.Partitions(appConfig.Topic)
        if err != nil {
            logs.Error("get partitions failed, err:", err)
            continue
        }

        for _, p := range partitions {
            partitionConsumer, err := consumer.ConsumePartition(appConfig.Topic, p, sarama.OffsetOldest)
            if err != nil {
                logs.Error("partitionConsumer err:", err)
                continue
            }
            wg.Add(1)
            go func(){
                for m := range partitionConsumer.Messages() {
                    logs.Debug("key: %s, text: %s, offset: %d\n", string(m.Key), string(m.Value), m.Offset)
                }
                wg.Done()
            }()
        }
        wg.Wait()
    }
    
    logs.Debug("Consumer success.")
}
