package elaeagnus

import (
	"fmt"
	// "github.com/go-elaeagnus/elaeagnus/httpcore"
	"github.com/go-elaeagnus/elaeagnus/tcpcore"
	"os"
)

func Run() {
	if len(os.Args) > 1 {
		if os.Args[1] == "stop" {
			fmt.Printf("stop command\n")
			tcpcore.StopCommand()
		}
		return
	}

	go Http(81)
	go tcpcore.Tcp(82, 0)
	tcpcore.Tcp(83, 1)
	fmt.Println("good bye elaeagnus")
}
