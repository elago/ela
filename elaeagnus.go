package ela

import (
// "fmt"
// "github.com/gogather/com/log"
// "reflect"
)

const (
	VERSION = "0.0.1"
)

var (
	config      = NewEmptyConfig()
	middlewares []interface{}
)

func Version() string {
	return VERSION
}

func Use(middleware interface{}) {
	// t := reflect.TypeOf(middleware)
	// if t.Kind() == reflect.Func {
	// 	result, err := injectFuc(middleware)
	// 	if err != nil {
	// 		log.Redf("injection failed: %s\n", err)
	// 	} else {
	// 		middleware = result[0]
	// 	}
	// }
	middlewares = append(middlewares, middleware)
}

func Run() {
	servHTTP(int(config.GetIntDefault("_", "port", 3000)))
}

func SetConfig(path string) {
	config.ReloadConfig(path)
}

func GetConfig() *Config {
	return config
}
