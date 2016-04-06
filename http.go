package ela

import (
	"fmt"
	"github.com/gogather/com/log"
	"net/http"
	"regexp"
	"runtime"
	"strings"
)

func servHTTP(port int) {
	log.Pinkf("[ela] Listen Port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &elaHandler{})
	if err != nil {
		fmt.Printf("HTTP Server Start Failed Port [%d]\n%s", port, err)
	}
}

// Http handler
type elaHandler struct{}

func (*elaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// package ctx
	ctx := Context{}
	ctx.w = NewResponseWriter(w)
	ctx.r = r
	ctx.Data = make(map[string]interface{})
	ctx.SetStatus(200)

	// parse router and call action
	path := r.URL.String()

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
				staticServ(path, ctx.w, r)
			} else {
				ctx.SetStatus(404)
				http.Error(ctx.w, "404, File Not Exist", 404)
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
				staticServ(rpath, ctx.w, r)
			} else {
				ctx.SetStatus(404)
				http.Error(ctx.w, "404, File Not Exist", 404)
			}
		}
	}

	f := uriMapping[path]
	if f != nil {
		function := f.(func(Context))

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

				content := "500 Server Internal Error!\n\n" + fmt.Sprintf("%s", r) + "\n\n" + stack
				log.Redln(r)
				log.Yellowln(stack)
				ctx.SetStatus(500)
				ctx.Write(content)
				responseLog(ctx)
			}
		}()

		// excute controller
		function(ctx)
	} else if errDefault != nil {
		// if static-alias does not exist, using default mode
		if staticExist(path) {
			staticServ(path, ctx.w, r)
		} else {
			ctx.SetStatus(404)
			http.Error(ctx.w, "404, File Not Exist", 404)
		}
	}

	// recording response log
	responseLog(ctx)
}
