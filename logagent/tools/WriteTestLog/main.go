package main

import (
    "encoding/json"
    "fmt"
    "time"
    "github.com/astaxie/beego/logs"
)

func main() {
    config := make(map[string]interface{})
    config["filename"] = "F:/Go/project/src/go_dev/LogCollectProject/logs/collectLog.log"
    config["level"] = logs.LevelDebug

    configStr, err := json.Marshal(config)
    if err != nil {
        fmt.Println("marshal failed, err:", err)
        return
    }
    logs.SetLogger(logs.AdapterFile, string(configStr))
    
    count := 0
    for {
        count += 1
        logs.Debug("this is a debug, count is %d", count)
        logs.Trace("this is a trace, count is %d", count)
        logs.Warn("this is a warn, count is %d", count)
        time.Sleep(time.Second)
    }
    
}
