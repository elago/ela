package ela

import (
	"fmt"
	"github.com/gogather/com"
	"strings"
	"errors"
)

func servError(ctx Context, err string, status int, useDefault bool) {
	f := getController(fmt.Sprintf("@%d", status))
	if f != nil && !useDefault{
		function := f.(func(Context, error))

		defer func() {
			if r := recover(); r != nil {
				ctx.SetHeader("Content-Type", "text/html")
				ctx.SetStatus(status)
				ctx.Write(err)
			}
		}()

		function(ctx, errors.New(err))
		ctx.SetStatus(status)
	} else {
		ctx.SetHeader("Content-Type", "text/html")
		ctx.SetStatus(status)
		ctx.Write(err)
	}
}

func fmtErrorHtml(content string) string {
	content = com.HTMLEncode(content)
	content = strings.Replace(content, "\n", "<br>", -1)
	return content
}
