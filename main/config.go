package main

import (
    "fmt"
    "errors"
    "go_dev/LogCollectProject/common"
    "github.com/astaxie/beego/config"
)

var (
    appConfig *Config = &Config{} //must init
)

type Config struct {
    //collect
    Collect_log_path string
    Topic string
    //logs
    Logs_log_path string
    Log_level string
    //kafka
    Kafka_addr string
}

func InitConfig(fileType, conFilePath string) (err error) {
    conf, err := config.NewConfig(fileType, conFilePath)
    if err != nil {
        err = common.ErrMsg("Config failed.", err)
        return
    }
    //load collect
    appConfig.Topic = conf.String("collect::topic")
    if len(appConfig.Topic) == 0 {
        err = errors.New("kafka topic is null.")
        return 
    }

    appConfig.Collect_log_path = conf.String("collect::log_path")
    if len(appConfig.Collect_log_path) == 0 {
        err = errors.New("Collect log path is null.")
        return 
    }
   
    //load logs
    appConfig.Logs_log_path = conf.String("logs::log_path")
    if (len(appConfig.Logs_log_path) == 0) {
        err = errors.New("logs path is null.")
        return 
    }

    appConfig.Log_level = conf.String("logs::log_level")
    if len(appConfig.Log_level) == 0 {
        appConfig.Log_level = "debug"
        fmt.Printf("Log level is null, use default %s.\n", appConfig.Log_level)
    }

    //load kafka
    appConfig.Kafka_addr = conf.String("kafka::server_addr")
    if len(appConfig.Kafka_addr) == 0 {
        err = errors.New("Kafka addr is null.")
        return 
    }

    fmt.Println("appConfig: ", appConfig)

    return
}