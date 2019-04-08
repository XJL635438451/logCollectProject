package main

import (
    "fmt"
    "github.com/astaxie/beego/logs"
)

func main() {
    conFilePath := "F:/Go/project/src/go_dev/LogCollectProject/conf/logAgent.conf"
    fileType := "ini"
    
    err := InitConfig(fileType, conFilePath)
    if err != nil {
        fmt.Println(err)
        return 
    }
    fmt.Println("config success")
    
    err = InitLog(appConfig.Logs_log_path, appConfig.Log_level)
    if err != nil {
        fmt.Println(err)
        return 
    }
    logs.Debug("init log success")

    tails, err := InitTailf(appConfig.Collect_log_path)
    if err != nil {
        logs.Error(err)
        return 
    }
    logs.Debug("init tailf success")

    err = InitKafka(appConfig.Kafka_addr)
    if err != nil {
        logs.Error(err)
        return 
    }
    logs.Debug("init kafka success")

    SendDataToKafka(tails)
}