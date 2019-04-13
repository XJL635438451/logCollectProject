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
	fmt.Printf("Load conf success.")

	err = initLogger()
	if err != nil {
		panic(err)
	}
	logs.Info("Load conf success, config:%v", appConfig)
    //init etcd
	collectConf, err := initEtcd()
	if err != nil {
		logs.Error("Init etcd failed, Error: %v", err)
		return
	}
	logs.Debug("initialize etcd success")
    //init tail
	err = tailf.InitTail(collectConf, appConfig.chanSize)
	if err != nil {
		logs.Error("init tail failed, err:%v", err)
		return
	}
	logs.Debug("initialize tailf success")
    //init kafka
	err = kafka.InitKafka(appConfig.kafkaAddr)
	if err != nil {
		logs.Error("Init tail failed, err:%v", err)
		return
	}
	logs.Debug("initialize kafka success")
    // init all success
	logs.Debug("Initialize all success")
	
	// go WriteTestLog()
	
	err = serverRun()
	if err != nil {
		logs.Error("serverRUn failed, err:%v", err)
		return
	}

	logs.Info("program exited")
}

// func WriteTestLog() {
// 	var count int
// 	for {
// 		count++
// 		logs.Debug("test for logger %d", count)
// 		time.Sleep(time.Millisecond * 1000)
// 	}
// }