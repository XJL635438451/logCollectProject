package tailf

import (
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

//一个日志收集的object
type TailObj struct {
	tail     *tail.Tail
	conf     CollectConf
	status   int  //normal or delete
	exitChan chan int  //exit
}

//要发给kafka的数据
type TextMsg struct {
	Msg   string
	Topic string
}

//总的objects管理
type TailObjMgr struct {
	tailObjs []*TailObj  //存放所有的日志收集的对象
	msgChan  chan *TextMsg //与kafka通信的管道，该管道的大小在配置文件中
	lock     sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

//需要注意：
//1. 新增日志收集路径obj
//   confs是从etcd中拿出来的(有可能已更新)，而此时tailObjMgr.tailObjs中还没有更新，
//   因此循环confs，如果confs中的LogPath在tailObjMgr.tailObjs中没找到，
//   说明是新增日志收集路径，创建并将新增的obj添加到tailObjMgr.tailObjs中。
//2. 删除日志收集路径obj
//   反过来，循环tailObjMgr.tailObjs，在confs中如果不存在前者的obj，说明已经被删除了，设置status，
//   更新tailObjMgr.tailObjs，退出被删除的obj日志收集程序。
func UpdateConfig(confs []CollectConf) (err error) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()
    //conf中的收集实例是否在tailObjMgr.tailObjs中，不在则是新的需要重新创建
	for _, oneConf := range confs {
		var isRunning = false
		//遍历所有的obj，看是否obj还在运行
		//因为有可能
		for _, obj := range tailObjMgr.tailObjs {
			if oneConf.LogPath == obj.conf.LogPath {
				isRunning = true
				break
			}
		}

		if isRunning {
			continue
		}

		createNewTask(oneConf)
	}
    //更新 tailObjs
	var tailObjs []*TailObj
	for _, obj := range tailObjMgr.tailObjs {
		obj.status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}

		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, obj)
	}

	tailObjMgr.tailObjs = tailObjs
	return
}

func createNewTask(conf CollectConf) {
	logs.Debug("Start to create new task.")
	obj := &TailObj{
		conf:     conf,
		exitChan: make(chan int, 1),
	}

	tails, errTail := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})

	if errTail != nil {
		logs.Error("Collect filename[%s] failed, err:%v", conf.LogPath, errTail)
		return
	}

	obj.tail = tails
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)
	logs.Debug("Successfully created new task and start to goroutine read from tail task..")
	go readFromTail(obj)
}

func InitTail(conf []CollectConf, chanSize int) (err error) {
	logs.Debug("Start to init tail.")
	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}

	if len(conf) == 0 {
		logs.Error("Invalid config for log collect, Conf: %v", conf)
		return
	}

	for _, v := range conf {
		createNewTask(v)
	}
    logs.Debug("Successfully initialized tail.")
	return
}

func readFromTail(tailObj *TailObj) {
	for {
		select {
		case line, ok := <-tailObj.tail.Lines:
			if !ok {
				logs.Warn("Tail file close reopen, filename:%s\n", tailObj.tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			textMsg := &TextMsg{
				Msg:   line.Text,
				Topic: tailObj.conf.Topic,
			}
            
			tailObjMgr.msgChan <- textMsg
		case <-tailObj.exitChan:
			logs.Warn("Tail obj will exited, Conf: %v", tailObj.conf)
			return
		}
	}
}
