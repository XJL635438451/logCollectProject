package main

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	"MyGitHubProject/logCollectProject/logagent/tailf"
)

var (
	appConfig *Config = &Config{}
)

type Config struct {
	//logs
	logLevel string
	logPath  string
    //kafka
	chanSize    int
	kafkaAddr   string
	collectConf []tailf.CollectConf //[] logpath, topic
    //etcd
	etcdAddr string
	etcdKey  string
	etcdDailTimeout int
}

//Init colect object config, add collect object by etcd
func loadCollectConf(conf config.Configer) (err error) {
    fmt.Println("Start to load collect config.")
	var cc tailf.CollectConf
	cc.LogPath = conf.String("collect::log_path")
	if len(cc.LogPath) == 0 {
		err = errors.New("Lnvalid collect::log_path.")
		return
	}

	cc.Topic = conf.String("collect::topic")
	if len(cc.Topic) == 0 {
		err = errors.New("Lnvalid collect::topic.")
		return
	}
	
	appConfig.collectConf = append(appConfig.collectConf, cc)
	fmt.Printf("Successfully loaded collect config.")
	return
}

//init logs config
func initLogsConf(conf config.Configer) (err error) {
	fmt.Println("Start to load logs config.")
	appConfig.logLevel = conf.String("logs::log_level")
	if len(appConfig.logLevel) == 0 {
		err = errors.New("Lnvalid logs::log_level.")
        return
	}

	appConfig.logPath = conf.String("logs::log_path")
	if len(appConfig.logPath) == 0 {
		err = errors.New("Lnvalid logs::log_path.")
        return
	}

	appConfig.chanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		err = fmt.Errorf("Lnvalid collect::chan_size. Error: %v", err)
        return
	}
	fmt.Printf("Successfully loaded logs config.")
	return
}

//init kafka config
func initKafkaConf(conf config.Configer) (err error) {
	fmt.Println("Start to load kafka config.")
	appConfig.kafkaAddr = conf.String("kafka::server_addr")
	if len(appConfig.kafkaAddr) == 0 {
		err = errors.New("Lnvalid kafka::server_addr.")
		return
	}
	fmt.Printf("Successfully loaded kafka config.")
    return
}

//init etcd config
func initEtcdConf(conf config.Configer) (err error) {
	fmt.Println("Start to load etcd config.")
	appConfig.etcdAddr = conf.String("etcd::addr")
	if len(appConfig.etcdAddr) == 0 {
		err = errors.New("Lnvalid etcd::addr.")
		return
	}

	appConfig.etcdKey = conf.String("etcd::configKey")
	if len(appConfig.etcdKey) == 0 {
		err = errors.New("Lnvalid etcd::configKey.")
		return
	}

	appConfig.etcdDailTimeout, err = conf.Int("etcd::etcdDailTimeout")
	if err != nil {
		err = fmt.Errorf("Lnvalid etcd::etcdDailTimeout. Error: %v", err)
		return
	}
	fmt.Printf("Successfully loaded etcd config.")
	return
}

func loadConf(confType, filename string) (err error) {
	fmt.Println("Start to load all config.")
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		err = fmt.Errorf("New config failed, Error: %v", err)
		return
	}

	err = initLogsConf(conf)
	if err != nil {
		err = fmt.Errorf("Failed to init logs config. Error: %v", err)
		return
	}

	err = initKafkaConf(conf)
	if err != nil {
		err = fmt.Errorf("Failed to init kafka config. Error: %v", err)
		return
	}
	
	err = initEtcdConf(conf)
	if err != nil {
		err = fmt.Errorf("Failed to init etcd config. Error: %v", err)
		return
	}
	
	err = loadCollectConf(conf)
	if err != nil {
		err = fmt.Errorf("Load collect conf failed, Error: %v", err)
		return
	}
	fmt.Printf("Successfully loaded all config, Config: %v\n", appConfig)
	return
}
