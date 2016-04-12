package ela

import (
	"fmt"
	"net/http"
)

// RequestContext
type Context struct {
	w         ResponseWriter
	r         *http.Request
	Data      map[string]interface{}
	status    int
	headerMap map[string]string
	uriParams map[string]string
}

func (ctx *Context) GetResponseWriter() ResponseWriter {
	return ctx.w
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.r
}

func (ctx *Context) GetMethod() string {
	return ctx.r.Method
}

func (ctx *Context) SetStatus(status int) {
	ctx.status = status
	ctx.w.SetStatus(status)
	// ctx.w.WriteHeader(ctx.status)
}

func (ctx *Context) GetStatus() int {
	return ctx.status
}

func (ctx *Context) SetHeader(key, value string) {
	if ctx.headerMap == nil {
		ctx.headerMap = make(map[string]string)
	}
	ctx.headerMap[key] = value
}

func (ctx *Context) GetCookie(key string) (*http.Cookie, error) {
	return ctx.r.Cookie(key)
}

func (ctx *Context) SetCookie(cookie *http.Cookie) {
	// log.Redln(cookie)
	http.SetCookie(ctx.w, cookie)
}

func (ctx *Context) Write(content string) (int, error) {
	header := ctx.w.Header()
	for k, v := range ctx.headerMap {
		header.Set(k, v)
		// ctx.r.Header.Set(k, v)
	}
	ctx.w.WriteHeader(ctx.status)
	return ctx.w.Write([]byte(content))
}

func (ctx *Context) serveTemplateWithStatus(templateFile string, status int, useDefaultError bool) {
	ctx.SetHeader("Content-Type", "text/html")

	// if in debug mode, reload templates
	if config.GetStringDefault("_", "mode", "dev") == "dev" {
		reloadTemplate()
	}

	t, err := parseFiles(templatesName...)

	if err != nil {
		content := "<h2>Server Internal Error!</h2>\n\n" + fmt.Sprintln(err)
		servError(*ctx, content, 500, useDefaultError)
	} else {

		if status == 404 {
			ctx.w.WriteHeader(404)
		} else if status == 500 {
			ctx.w.WriteHeader(500)
		}

		err = t.ExecuteTemplate(ctx.w, templateFile, ctx.Data)

		if err != nil {
			content := "<h2>Server Internal Error!</h2>\n\n" + fmt.Sprintf("<pre>%s</pre>", fmtErrorHtml(err.Error()))
			servError(*ctx, content, 500, useDefaultError)
		}
	}

}

func (ctx *Context) ServeTemplate(templateFile string) {
	ctx.serveTemplateWithStatus(templateFile, 200, false)
}

func (ctx *Context) ServeError(status int, templateFile string) {
	ctx.serveTemplateWithStatus(templateFile, status, true)
}

func (ctx *Context) setURIParam(params map[string]string) {
	ctx.uriParams = params
}

func (ctx *Context) GetURIParam(key string) (string, error) {
	if ctx.uriParams == nil {
		return "", fmt.Errorf("%s", "does not exist uri params")
	} else {
		// fmt.Printf("%s", ctx.uriParams)
		return ctx.uriParams[key], nil
	}
}

func (ctx *Context) GetURIParamDefault(key string, defaultValue string) {

}
