package main

import (
    "encoding/json"
    "strings"
    "github.com/astaxie/beego/logs"
    "go_dev/LogCollectProject/common"
)

func convertLogLevel(level string) int {
    level = strings.ToLower(level)

    switch (level) {
    case "debug":
        return logs.LevelDebug
    case "warn":
        return logs.LevelWarn
    case "info":
        return logs.LevelInfo
    case "trace":
        return logs.LevelTrace
    }

    return  logs.LevelDebug
}

func InitLog(filePath, logLevel string) (err error) {
    config := make(map[string]interface{})
    config["filename"] = filePath
    config["level"] = convertLogLevel(logLevel)

    configStr, err := json.Marshal(config)
    if err != nil {
        err = common.ErrMsg("Marshal failed, err:", err)
        return
    }

    logs.SetLogger(logs.AdapterFile, string(configStr))
    return
}
