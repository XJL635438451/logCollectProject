package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/astaxie/beego/logs"
	etcd_client "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"MyGitHubProject/logCollectProject/logagent/tailf"
)

type EtcdClient struct {
	client *etcd_client.Client
	keys   []string //存放etcd中所有的key
}

var (
	etcdClient *EtcdClient
)

//注意：这块的key只是真正etcd的key前缀，后面会和每台机器的ip拼接起来，作为不同机器etcd的key
//1. 如果etcd中有etcdKey则取出对应的value，程序结束会返回
//2. 没有该etcdKey则watch
func initEtcd() (collectConf []tailf.CollectConf, err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{appConfig.etcdAddr},
		DialTimeout: time.Duration(appConfig.etcdDailTimeout) * time.Second,
	})
	if err != nil {
		err = fmt.Errorf("Connect etcd failed, Error: %v", err)
		return
	}

	etcdClient = &EtcdClient{
		client: cli,
	}
    var key string = appConfig.etcdKey
	if strings.HasSuffix(key, "/") == false {
		key = appConfig.etcdKey + "/"
	}
    //localIPArray只有当前机器的ip
	for _, ip := range localIPArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		etcdClient.keys = append(etcdClient.keys, etcdKey)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("Client get from etcd failed, Error: %v", err)
			continue
		}
		cancel()

		logs.Debug("Resp from etcd: %v", resp.Kvs)
		//resp.Kvs --> [key:"/logagent/conf/" create_revision:2 mod_revision:14 version:13 value:"sample_value" ]
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				//1. 在收集多台机器的日志时用的是ip来区分不同机器
				//2. 一台机器可能会收集多个日志
				// collectConf = [{"path":"D:/project/nginx/logs/access2.log","topic":"nginx_log"},
                // {"path":"D:/project/nginx/logs/error2.log","topic":"nginx_log_err"}]
				err = json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("unmarshal failed, Error: %v", err)
					continue
				}
				logs.Info("Log collect config is %v", collectConf)
			}
		}
	}

	initEtcdWatcher()
	return
}

func initEtcdWatcher() {
	//如果key不存在也是可以watch，等该key发生变化时则watch就可以察觉
	for _, key := range etcdClient.keys {
		go watchKey(key)
	}
}

func watchKey(key string) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{appConfig.etcdAddr},
		DialTimeout:  time.Duration(appConfig.etcdDailTimeout) * time.Second,
	})
	if err != nil {
		logs.Error("Connect etcd failed, Error: %v", err)
		return
	}

	logs.Debug("Begin watch key: %s", key)
	for {
		//rch 是一个 WatchChan 类型的管道 --> (WatchChan <-chan WatchResponse)
		rch := cli.Watch(context.Background(), key)
		var collectConf []tailf.CollectConf
		for wresp := range rch {
			for _, ev := range wresp.Events {
				//1. 如果本来就没有该key，那也就不涉及删除
				//2. 如果key在etcd中，当key被删除了，后面会更新objs
				if ev.Type == mvccpb.DELETE {
					logs.Warn("Key[%s] 's config is deleted", key)
					continue
				}
                //增加key或者修改key（也就相当于删除了原来的key）
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", key, ev.Kv.Value, err)
						continue
					}
					logs.Debug("Get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					logs.Debug("Get config from etcd succ, %v", collectConf)

					tailf.UpdateConfig(collectConf)
				}
			}
		}

	}
}
