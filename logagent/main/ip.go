package main

import (
	"fmt"
	"net"
)

var (
	localIPArray []string
)

func init() {
	netInterfaces, err := net.Interfaces()
    if err != nil {
        panic(fmt.Sprintf("Get local ip failed, %v", err))
    }
    // 判断net.FlagUp标志进行确认，排除掉无用的网卡。
    for i := 0; i < len(netInterfaces); i++ {
        if (netInterfaces[i].Flags & net.FlagUp) != 0 {
            addrs, _ := netInterfaces[i].Addrs()
 
            for _, address := range addrs {
                if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                    if ipnet.IP.To4() != nil {
						localIPArray = append(localIPArray, ipnet.IP.String())
                        return
                    }
                }
            }
        }
	}
	// fmt.Println(localIPArray)
    return
}

