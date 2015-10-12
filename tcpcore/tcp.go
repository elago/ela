package tcpcore

import (
	"fmt"
	// "io"
	"net"
	"os"
)

func Tcp(port int, maxClientNum int) {
	i := 0
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	handleError(err, nil)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			continue
		}
		i += 1
		if maxClientNum > 0 && i > maxClientNum {
			fmt.Println("reached max client limit, server stoped.")
			return
		}
		go handleConnection(conn, i)
	}

}

func handleConnection(tcpConn net.Conn, cid int) {
	fmt.Println("client [" + fmt.Sprintf("%s", cid) + "] connected")
}

func handleError(err error, tcpConn net.Conn) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Client error: %s\n", err.Error())
		if tcpConn != nil {
			tcpConn.Close()
		}
	}
}
