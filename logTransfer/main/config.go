package main

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

type LogConfig struct {
	LogPath    string
	LogLevel   string

	KafkaAddr  string
	KafkaTopic string

	ESAddr     string
}

var (
	logConfig *LogConfig = &LogConfig{}
)

func initConfig(confType string, filename string) (err error) {
	fmt.Println("Start to init config.")
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		err = fmt.Errorf("New config failed, err: %v", err)
		return
	}

	//init logs config
	logConfig.LogLevel = conf.String("logs::log_level")
	if len(logConfig.LogLevel) == 0 {
		logConfig.LogLevel = "debug"
	}

	logConfig.LogPath = conf.String("logs::log_path")
	if len(logConfig.LogPath) == 0 {
		logConfig.LogPath = "./logs"
	}

	//init kafka config
	logConfig.KafkaAddr = conf.String("kafka::server_addr")
	if len(logConfig.KafkaAddr) == 0 {
		err = fmt.Errorf("Invalid kafka addr")
		return
	}

	logConfig.KafkaTopic = conf.String("kafka::topic")
	if len(logConfig.ESAddr) == 0 {
		err = fmt.Errorf("Invalid es addr")
		return
	}

	//init es config
	logConfig.ESAddr = conf.String("es::addr")
	if len(logConfig.ESAddr) == 0 {
		err = fmt.Errorf("Invalid es addr")
		return
	}

	fmt.Printf("Successful initialized config, logConfig: %v\n", logConfig)
	return
}
