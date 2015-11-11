package debug

import "fmt"

const (
	DEBUG = true
)

func Log(content string) {
	if !DEBUG {
		return
	} else {
		fmt.Println(content)
	}
}
