package elaeagnus

import (
	"fmt"
	"github.com/go-elaeagnus/elaeagnus/httpcore"
	"github.com/go-elaeagnus/elaeagnus/tcpcore"
	"net"
	"os"
	// "strings"
)

const (
	addr = "127.0.0.1:83"
)

var manager net.Conn
var connectValid bool

func init() {
	connectValid = true
}

func connectManager() {
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

func stopCommand() {
	connectManager()
	if connectValid {
		_, err := manager.Write([]byte("stop"))
		if err != nil {
			fmt.Println("Send Failed: " + err.Error())
		}
	} else {
		fmt.Println("命令发送失败")
	}
}

func Run() {
	if len(os.Args) > 1 {
		if os.Args[1] == "stop" {
			fmt.Printf("stop command\n")
			stopCommand()
		}
		return
	}

	go httpcore.Http(81)
	go tcpcore.Tcp(82, 0)
	tcpcore.Tcp(83, 1)
	fmt.Println("good bye elaeagnus")
}
