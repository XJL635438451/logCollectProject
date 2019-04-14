package main

import (
	"github.com/astaxie/beego/logs"
)

func main() {
	fileName := "./conf/log_transfer.conf"
	confType := "ini"
	err := initConfig(confType, fileName)
	if err != nil {
		panic(err)
	}

	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
	}
	
	err = initKafka(logConfig.KafkaAddr, logConfig.KafkaTopic)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}

	err = initES(logConfig.ESAddr)
	if err != nil {
		logs.Error("Init es failed, err:%v", err)
		return
	}

	err = run()
	if err != nil {
		logs.Error("Run  failed, err:%v", err)
		return
	}

	logs.Warn("Main program exception, log_transfer exited.")
}
