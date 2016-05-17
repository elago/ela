package ela

import (
	"errors"
	"fmt"
	"github.com/gogather/com"
	"strings"
)

func servError(ctx Context, err string, status int, useDefault bool) {
	controller := getController(fmt.Sprintf("@%d", status))

	if controller == nil {
		servError(ctx, "<h2>404, File Not Exist</h2>", 404, false)
		return
	}

	routerElement := controller.(uriMode)
	f := routerElement.fun

	if f != nil && !useDefault {
		functions := f

		defer func() {

			if r := recover(); r != nil {
				ctx.SetHeader("Content-Type", "text/html")
				ctx.SetStatus(status)
				ctx.Write(err)
			}
		}()

		ctx.SetStatus(status)

		// just get and execute first controller
		if len(functions) >= 1 {
			function := functions[0]
			Use(errors.New(err))
			injectFuc(function)
		}

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
