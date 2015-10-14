package tcpcore

import (
	"fmt"
	"net"
)

const (
	addr = "127.0.0.1:83"
)

var manager net.Conn
var connectValid bool

func init() {
	connectValid = true
}

func ConnectManager() {
	var err error
	manager, err = net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		connectValid = false
		handleConnectFailed()
		return
	} else {
		if manager == nil {
			fmt.Println("链接失败")
			connectValid = false
			handleConnectFailed()
		}
		connectValid = true
	}
	fmt.Println("已连接服务器")
}

func handleConnectFailed() {
	if manager != nil {
		manager.Close()
	}
}
