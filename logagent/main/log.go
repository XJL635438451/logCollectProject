package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"github.com/astaxie/beego/logs"
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


func initLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = appConfig.logPath
	config["level"] = convertLogLevel(appConfig.logLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		err = fmt.Errorf("InitLogger failed, marshal err: %v", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}