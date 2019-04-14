package main

import (
	"fmt"

	elastic "gopkg.in/olivere/elastic.v2"
	"github.com/astaxie/beego/logs"
)

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var (
	esClient *elastic.Client
)

func initES(addr string) (err error) {
    logs.Debug("Start to init es.")
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if err != nil {
		fmt.Println("connect es error", err)
		return
	}
	esClient = client
    logs.Debug("Successful initialized es.")
	return
}

func sendToES(topic string, data []byte) (err error) {
	msg := &LogMessage{}
	msg.Topic = topic
	msg.Message = string(data)

	_, err = esClient.Index().
		Index(topic).
		Type(topic).
		BodyJson(msg).
		Do()
	if err != nil {
		panic(err)
	}
	return
}
