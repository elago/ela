package ela

import (
	// "errors"
	"fmt"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	// "reflect"
	"strings"
)

func servError(ctx *Context, err string, status int, useDefault bool) {
	inj := newInjection()

	controller := getController(fmt.Sprintf("@%d", status))

	if controller == nil {
		ctx.SetHeader("Content-Type", "text/html")
		ctx.SetStatus(status)
		ctx.Write(err)
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
			inj.headMiddleware(ctx)
			inj.appendMiddleware(err)
			_, err := inj.injectFuc(function)
			if err != nil {
				log.Redf("injection failed: %s\n", err)
			}

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
