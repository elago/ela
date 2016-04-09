package ela

import (
	"fmt"
	"net/http"
	"github.com/gogather/com/log"
)

func servError(ctx Context, err string, status int) {
	f := uriMapping[fmt.Sprintf("@%d", status)]
	if f != nil {
		function := f.(func(Context))

		defer func() {
			if r := recover(); r != nil {
				log.Pinkln("===")
				http.Error(ctx.w, err, 500)
			}
		}()

		function(ctx)
		ctx.SetStatus(500) // TODO: response writer should return 500 status
	} else {
		http.Error(ctx.w, err, 500)
	}
}
