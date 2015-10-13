package httpcore

import (
	"fmt"
	"io"
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
	io.WriteString(w, "connected.")
}
