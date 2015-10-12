package elaeagnus

import (
	"fmt"
	"github.com/go-elaeagnus/elaeagnus/httpcore"
	"github.com/go-elaeagnus/elaeagnus/tcpcore"
)

func Run() {
	go httpcore.Http(81)
	go tcpcore.Tcp(82, 0)
	tcpcore.Tcp(83, 1)
	fmt.Println("good bye elaeagnus")
}
