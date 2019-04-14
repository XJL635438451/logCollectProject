package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"MyGitHubProject/logCollectProject/logagent/kafka"
	"MyGitHubProject/logCollectProject/logagent/tailf"
)

func main() {
	fileName := "./conf/logagent.conf"
	confType := "ini" 
	err := loadConf(confType, fileName)
	if err != nil {
		err = fmt.Errorf("Load conf falied, Error: %v", err)
		panic(err)
	}
	
	err = initLogger()
	if err != nil {
		panic(err)
	}

	collectConf, err := initEtcd()
	if err != nil {
		logs.Error("Init etcd failed, Error: %v", err)
		return
	}

	err = tailf.InitTail(collectConf, appConfig.chanSize)
	if err != nil {
		logs.Error("init tail failed, err:%v", err)
		return
	}

	err = kafka.InitKafka(appConfig.kafkaAddr)
	if err != nil {
		logs.Error("Init tail failed, err:%v", err)
		return
	}
	
	err = serverRun()
	if err != nil {
		logs.Error("serverRUn failed, err:%v", err)
		return
	}

	logs.Info("program exited")
}