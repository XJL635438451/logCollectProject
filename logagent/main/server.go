package main

import (
	"fmt"
	"MyGitHubProject/logCollectProject/logagent/kafka"
	"MyGitHubProject/logCollectProject/logagent/tailf"
	"time"

	"github.com/astaxie/beego/logs"
)

func serverRun() (err error) {
	for {
		msg := tailf.GetOneLine()
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("send to kafka failed, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	fmt.Printf("Read msg:%s, topic:%s\n", msg.Msg, msg.Topic)
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
