package ela

import (
	"fmt"
	"github.com/gogather/com"
	// "github.com/gogather/com/log"
	"net/http"
	"strings"
)

func servError(ctx Context, err string, status int, useDefault bool) {
	// err = fmtErrorHtml(err)
	f := uriMapping[fmt.Sprintf("@%d", status)]
	if f != nil && !useDefault{
		function := f.(func(Context))

		defer func() {
			if r := recover(); r != nil {
				ctx.SetHeader("Content-Type", "text/html")
				http.Error(ctx.w, err, status)
			}
		}()

		function(ctx)
		ctx.SetStatus(status)
	} else {
		ctx.SetHeader("Content-Type", "text/html")
		http.Error(ctx.w, err, status)
	}
}

func fmtErrorHtml(content string) string {
	content = com.HTMLEncode(content)
	content = strings.Replace(content, "\n", "<br>", -1)
	return content
}
