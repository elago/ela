package httpcore

import (
	// "github.com/gogather/com/log"
	"fmt"
	"net/http"
)

func Http(port int) {
	http.HandleFunc("/", HandleJsonRpc)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Tcp Server Start Failed Port [%d]\n%s", port, err)
	}
}

func HandleJsonRpc(w http.ResponseWriter, r *http.Request) {

}
