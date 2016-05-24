package ela

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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
func (ctx *Context) GetCookie(key string) string {
	cookie, err := ctx.r.Cookie(key)
	if err != nil {
		return ""
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val
}

// set cookie
// args: name, value, max age, path, domain, secure, http only, expires
func (ctx *Context) SetCookie(name string, value string, others ...interface{}) {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = url.QueryEscape(value)

	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			cookie.MaxAge = v
		case int64:
			cookie.MaxAge = int(v)
		case int32:
			cookie.MaxAge = int(v)
		}
	}

	cookie.Path = "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	}

	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	if len(others) > 3 {
		switch v := others[3].(type) {
		case bool:
			cookie.Secure = v
		default:
			if others[3] != nil {
				cookie.Secure = true
			}
		}
	}

	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			cookie.HttpOnly = true
		}
	}

	if len(others) > 5 {
		if v, ok := others[5].(time.Time); ok {
			cookie.Expires = v
			cookie.RawExpires = v.Format(time.UnixDate)
		}
	}

	ctx.w.Header().Add("Set-Cookie", cookie.String())
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

// 10M memory for multipart form parsing
var MaxMemory = int64(1024 * 1024 * 10)

func (ctx *Context) parseForm() {
	if ctx.r.Form != nil {
		return
	}

	contentType := ctx.r.Header.Get("Content-Type")
	if (ctx.r.Method == "POST" || ctx.r.Method == "PUT") &&
		len(contentType) > 0 && strings.Contains(contentType, "multipart/form-data") {
		ctx.r.ParseMultipartForm(MaxMemory)
	} else {
		ctx.r.ParseForm()
	}
}

func (ctx *Context) GetParam(name string) string {
	ctx.parseForm()
	return ctx.r.Form.Get(name)
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
