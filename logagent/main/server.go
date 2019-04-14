package main

import (
	"MyGitHubProject/logCollectProject/logagent/kafka"
	"MyGitHubProject/logCollectProject/logagent/tailf"
	"time"

	"github.com/astaxie/beego/logs"
)

func serverRun() (err error) {
    logs.Debug("Start to run...")
	for {
		msg := tailf.GetOneLine()
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("Send message to kafka failed, Error: %v", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	logs.Debug("Start to send message to kafka...")
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
