package ela

import (
	"fmt"
	"github.com/codegangsta/inject"
	"github.com/gogather/com/log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

const (
	VERSION = "0.0.1"
)

func Version() string {
	return VERSION
}

var (
	config          = NewEmptyConfig()
	middlewares     = make(map[reflect.Type]interface{})
	middlewaresName []reflect.Type
)

// Http handler
type Elaeagnus struct {
}

func Web() *Elaeagnus {
	return &Elaeagnus{}
}

// define a serials controller as main controller for an uri pattern
func (ela *Elaeagnus) Router(uri string, f ...interface{}) {
	if strings.HasPrefix(uri, "@") {
		panic("@ should not be prefix of uri")
	} else if !strings.HasPrefix(uri, "/") {
		panic("uri should begin with /")
	} else {
		router(uri, true, true, f...)
	}
}

func (ela *Elaeagnus) execInjection() {
	var mids []reflect.Type
	var midw interface{}

	for i := 0; i < len(middlewaresName); i++ {
		mid := middlewaresName[i]
		if mid.Kind() == reflect.Func {
			// get the finnal object as middleware object
			midw = ela.parseInjectionObject(middlewares[mid])
			ela.Use(midw)

			if midw != nil {
				mids = append(mids, mid)
			}
		} else {
			ela.Use(middlewares[mid])
			mids = append(mids, mid)
		}

	}

	middlewaresName = mids
}

func (ela *Elaeagnus) parseInjectionObject(mid interface{}) interface{} {
	t := reflect.TypeOf(mid)
	if t.Kind() == reflect.Func {
		result, err := injectFuc(mid, middlewaresName...)
		if err != nil {
			log.Redf("injection failed: %s :L51\n", err)
			return nil
		} else {
			mid = result[0]
			return ela.parseInjectionObject(mid)
		}
	} else {
		return mid
	}
}

func (ela *Elaeagnus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// package ctx
	ctx := newContext(w, r)

	// add ctx in the head of middleware
	var mids []reflect.Type
	mids = append(mids, reflect.TypeOf(&ctx))
	middlewaresName = append(mids, middlewaresName...)
	ela.Use(&ctx)

	// execute injection
	ela.execInjection()

	// parse router and call action
	path := parseURI(r.URL.String())

	// recording request log
	requestLog(ctx)

	// add special static files
	addSpecialStatic("/favicon.ico")
	addSpecialStatic("/robots.txt")

	// deal with special static files
	for i := 0; i < len(specialStatic); i++ {
		special := specialStatic[i]
		if path == special {
			if staticExist(path) {
				staticServ(path, ctx)
			} else {
				servError(ctx, "<h2>404, File Not Exist</h2>", 404, false)
			}

			// recording response log
			responseLog(ctx)
			return
		}
	}

	staticAlias, errDefault := config.GetString("static", "alias")
	// if static-alias exist, using alias uri mode
	if errDefault == nil {
		if strings.HasPrefix(path, "/"+staticAlias) {
			reg := regexp.MustCompile(`^/` + staticAlias)
			rpath := reg.ReplaceAllString(path, "")
			path = reg.ReplaceAllString(path, "/"+staticDirectory)

			if staticExist(rpath) {
				staticServ(rpath, ctx)
			} else {
				servError(ctx, "<h2>404, File Not Exist</h2>", 404, false)
			}

			// recording response log
			responseLog(ctx)
			return
		}

	}

	ela.servController(path, ctx)

	// recording response log
	responseLog(ctx)
}

func (ela *Elaeagnus) servHTTP(port int) {
	log.Pinkf("[ela] Listen Port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &Elaeagnus{})
	if err != nil {
		fmt.Printf("HTTP Server Start Failed Port [%d]\n%s", port, err)
	}
}

func (ela *Elaeagnus) servController(path string, ctx Context) {
	controller := getController(path)
	if controller == nil {
		servError(ctx, "<h2>404, File Not Exist</h2>", 404, false)
		return
	}

	routerElement := controller.(uriMode)
	f := routerElement.fun
	params := routerElement.argsMap

	ctx.setURIParam(params)
	if f != nil {

		functions := f

		defer func() {
			if r := recover(); r != nil {

				var stack string
				for i := 1; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
				}

				content := "<h2>500 Server Internal Error!</h2>\n\n" + fmt.Sprintf("%s", r) + "\n\n" + "<pre>" + stack + "</pre>"
				log.Redln(r)
				log.Yellowln(stack)

				servError(ctx, content, 500, false)
				return
			}
		}()

		// execute before controllers
		if beforeController != nil && routerElement.withBefore {
			_, err := injectFuc(beforeController, middlewaresName...)
			if err != nil {
				log.Redf("injection failed: %s\n", err)
			}
		}

		// execute controllers
		for i := 0; i < len(functions); i++ {
			if !ctx.GetResponseWriter().HasFlushed() {
				function := functions[i]
				_, err := injectFuc(function, middlewaresName...)
				if err != nil {
					log.Redf("injection failed: %s\n", err)
				}
			}
		}

		// execute after controllers
		if afterController != nil && routerElement.withBefore {
			_, err := injectFuc(afterController, middlewaresName...)
			if err != nil {
				log.Redf("injection failed: %s\n", err)
			}
		}

	} else {
		// if static-alias does not exist, using default mode
		if staticExist(path) {
			staticServ(path, ctx)
		} else {
			servError(ctx, "<h2>404, File Not Exist</h2>", 404, false)
		}
	}
}

func (ela *Elaeagnus) Run() {
	ela.servHTTP(int(config.GetIntDefault("_", "port", 3000)))
}

func (ela *Elaeagnus) Use(middleware interface{}) {
	addMiddleware(middleware)
}

func SetConfig(path string) {
	config.ReloadConfig(path)
}

func GetConfig() *Config {
	return config
}

func addMiddleware(middleware interface{}) {
	t := reflect.TypeOf(middleware)
	middlewares[t] = middleware
	for i := 0; i < len(middlewaresName); i++ {
		element := middlewaresName[i]
		if element == t {
			return
		}
	}
	middlewaresName = append(middlewaresName, t)
}

func injectFuc(fun interface{}, mid ...reflect.Type) ([]reflect.Value, error) {
	inj := inject.New()
	for i := 0; i < len(mid); i++ {
		element := mid[i]
		if element.Kind() != reflect.Func {
			inj.Map(middlewares[element])
		}
	}
	return inj.Invoke(fun)
}
