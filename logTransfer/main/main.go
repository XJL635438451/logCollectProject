package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

func main() {
	fileName := "./conf/log_transfer.conf"
	confType := "ini" 
	err := initConfig(confType, fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Init config success.")

	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
	}
	logs.Debug("Init logger success.")

	err = initKafka(logConfig.KafkaAddr, logConfig.KafkaTopic)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}
	logs.Debug("Init kafka succ")

	err = initES(logConfig.ESAddr)
	if err != nil {
		logs.Error("init es failed, err:%v", err)
		return
	}
	logs.Debug("Init es client succ")

	err = run()
	if err != nil {
		logs.Error("run  failed, err:%v", err)
		return
	}

	logs.Warn("warning, log_transfer is exited")
}
