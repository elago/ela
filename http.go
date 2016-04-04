package ela

import (
	"fmt"
	"github.com/gogather/com/log"
	"net/http"
	"runtime"
)

func ServHttp(port int) {
	log.Pinkf("[ela] Listen Port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &ElaHandler{})
	if err != nil {
		fmt.Printf("HTTP Server Start Failed Port [%d]\n%s", port, err)
	}
}

// Http handler
type ElaHandler struct{}

func (*ElaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// package ctx
	ctx := Context{}
	ctx.w = NewResponseWriter(w)
	ctx.r = r
	ctx.Data = make(map[string]interface{})
	ctx.SetStatus(200)
	// parse router and call action
	path := r.URL.String()

	requestLog(ctx)

	f := URIMapping[path]
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

		function(ctx)
	} else {
		// show uri
		if StaticExist(path) {
			StaticServ(path, ctx.w, r)
		} else {
			ctx.SetStatus(404)
			http.Error(ctx.w, "404, File Not Exist", 404)
		}
	}

	responseLog(ctx)
}
