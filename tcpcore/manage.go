package tcpcore

import (
	"fmt"
)

func StopCommand() {
	ConnectManager()
	if connectValid {
		_, err := manager.Write([]byte("stop"))
		if err != nil {
			fmt.Println("Send Failed: " + err.Error())
		}
	} else {
		fmt.Println("命令发送失败")
	}
}
