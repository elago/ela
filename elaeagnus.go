package ela

import (
// "fmt"
)

var (
	config = NewConfig("conf/app.ini")
)

func SetConfig(path string) {
	config.ReloadConfig(path)
}

func Run() {
	ServHttp(int(config.GetIntDefault("_", "port", 3000)))
}
