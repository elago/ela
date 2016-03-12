package debug

import (
	"fmt"
	"github.com/gogather/com/log"
)

const (
	DEBUG = true
)

var STATUS map[int]string

func init() {
	STATUS = map[int]string{
		200: "ok",
		404: "page could not find",
		500: "server internal error",
		302: "not modified",
	}
}

// type struct{}

func requestPrint(status int, content string) {
	if status == 404 {
		log.Pinkln(content)
	} else if status == 500 {
		log.Redln(content)
	} else {
		log.Println(content)
	}
}

func RequestLog(status int, reqType string, method string, path string) {

	statusVal, ok := STATUS[status]
	if !ok {
		statusVal = ""
	}

	requestPrint(status, "[Ela] ["+reqType+"/"+method+"] "+path+" "+fmt.Sprintf("%d", status)+" "+statusVal)
}

func StartRequestLog(msg string) {
	log.Println("[Ela] " + msg)
}
