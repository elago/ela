package ela

import (
	"fmt"
	"github.com/elago/ela/debug"
	"net/http"
)

func Http(port int) {
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

	debug.StartRequestLog("starting get page " + path)

	f := URIMapping[path]
	if f != nil {
		function := f.(func(Context))
		function(ctx)
		debug.RequestLog(ctx.GetStatus(), "action", r.Method, path)
	} else {
		// show uri
		if StaticExist(path) {
			StaticServ(path, ctx.w, r)
		} else {
			ctx.SetStatus(404)
			http.Error(ctx.w, "404, File Not Exist", 404)
			debug.RequestLog(ctx.GetStatus(), "404", r.Method, path)
		}
	}
}
