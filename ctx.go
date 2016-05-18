package ela

import (
	"fmt"
	"net/http"
)

type Locale interface {
	Language() string
	Tr(string) string
}

// RequestContext
type Context struct {
	w         ResponseWriter
	r         *http.Request
	Data      map[string]interface{}
	status    int
	headerMap map[string]string
	uriParams map[string]string
	Locale
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

// url redirection
func (ctx *Context) Redirect(url string) {
	http.Redirect(ctx.w, ctx.r, url, 302)
}

// get cookie
func (ctx *Context) GetCookie(key string) (*http.Cookie, error) {
	return ctx.r.Cookie(key)
}

// set cookie
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.w, cookie)
}

// write and flush response content
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

// serve and parse a template
func (ctx *Context) ServeTemplate(templateFile string) {
	ctx.serveTemplateWithStatus(templateFile, 200, false)
}

// serve and parse a template for an error, used for error controller
func (ctx *Context) ServeError(status int, templateFile string) {
	ctx.serveTemplateWithStatus(templateFile, status, true)
}

func (ctx *Context) setURIParam(params map[string]string) {
	ctx.uriParams = params
}

// get uri params when defined router with uri params mode
func (ctx *Context) GetURIParam(key string) (string, error) {
	if ctx.uriParams == nil {
		return "", fmt.Errorf("%s", "does not exist uri params")
	} else {
		return ctx.uriParams[key], nil
	}
}

// get uri params with default value
func (ctx *Context) GetURIParamDefault(key string, defaultValue string) string {
	value, err := ctx.GetURIParam(key)
	if err != nil {
		return defaultValue
	} else {
		return value
	}
}
