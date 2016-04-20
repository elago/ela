package ela

import (
// "fmt"
)

const (
	VERSION = "0.0.1"
)

var (
	config = NewEmptyConfig()
)

func Version() string {
	return VERSION
}

func SetConfig(path string) {
	config.ReloadConfig(path)
}

func GetConfig() *Config {
	return &config
}

func Run() {
	servHTTP(int(config.GetIntDefault("_", "port", 3000)))
}
