package elaeagnus

import (
	"debug"
	"net/http"
)

func Http(port int) {
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), &ElaHandler{})
	if err != nil {
		fmt.Printf("Tcp Server Start Failed Port [%d]\n%s", port, err)
	}
}

type ElaHandler struct{}

func (*ElaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()
	debug.Log(path)
	f := URIMapping[path]
	if f != nil {
		function := f.(func())
		function()
	}
}
