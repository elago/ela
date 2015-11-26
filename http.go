package ela

import (
	"fmt"
	"github.com/elago/ela/debug"
	"net/http"
)

func Http(port int) {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &ElaHandler{})
	if err != nil {
		fmt.Printf("Tcp Server Start Failed Port [%d]\n%s", port, err)
	}
}

// RequestContext
type RequestContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (this *RequestContext) Write(content string) (int, error) {
	return this.w.Write([]byte(content))
}

// Http handler
type ElaHandler struct{}

func (*ElaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// package ctx
	ctx := RequestContext{}
	ctx.w = w
	ctx.r = r
	// parse router and call action
	path := r.URL.String()
	debug.Log(path)
	f := URIMapping[path]
	if f != nil {
		function := f.(func(RequestContext))
		function(ctx)
	}
}
