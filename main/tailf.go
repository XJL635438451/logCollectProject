package main

import (
    "github.com/hpcloud/tail"
    "go_dev/LogCollectProject/common"
)

func InitTailf(filePath string) (tails *tail.Tail,err error) {
    //filename := ".\\my.log"
    tails, err = tail.TailFile(filePath, tail.Config{
        ReOpen:    true,
        Follow:    true,
        //Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
        MustExist: false,   
        Poll:      true,
    })
    if err != nil {
        err = common.ErrMsg("tail file failed.", err)
        return
    }
    return
}
