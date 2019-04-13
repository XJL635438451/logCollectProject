package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

func run() (err error) {
	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic)
	if err != nil {
		logs.Error("Failed to get the list of partitions: ", err)
		return
	}

	for partition := range partitionList {
		pc, errRet := kafkaClient.client.ConsumePartition(kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			err = errRet
			logs.Error("Failed to consumer for partition %d, Error: %s\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		
		go func(pc sarama.PartitionConsumer) {
			kafkaClient.wg.Add(1)
			for msg := range pc.Messages() {
				logs.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				//fmt.Println()
				err = sendToES(kafkaClient.topic, msg.Value)
				if err != nil {
					logs.Warn("Send to es failed, topic:%v, value:%v, err:%v", kafkaClient.topic, msg.Value, err)
				}
			}
			kafkaClient.wg.Done()
		}(pc)
	}

	kafkaClient.wg.Wait()
	return
}
