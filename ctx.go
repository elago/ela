package ela

import (
	"fmt"
	// "github.com/gogather/com/log"
	"html/template"
	"net/http"
)

// RequestContext
type Context struct {
	w         ResponseWriter
	r         *http.Request
	Data      map[string]interface{}
	status    int
	headerMap map[string]string
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
	for k, v := range ctx.headerMap {
		ctx.r.Header.Set(k, v)
	}
	return ctx.w.Write([]byte(content))
}

func (ctx *Context) ServeTemplate(templateFile string) {
	ctx.SetHeader("Content-Type", "text/html")

	// if in debug mode, reload templates
	if config.GetStringDefault("_", "mode", "dev") == "dev" {
		reloadTemplate()
	}

	t, err := ctx.parseFiles(templatesName...)

	if err != nil {
		content := "Server Internal Error!\n\n" + fmt.Sprintln(err)
		ctx.SetStatus(500)
		ctx.Write(content)
	} else {
		err = t.ExecuteTemplate(ctx.w, templateFile, ctx.Data)
		content := "Server Internal Error!\n\n" + fmt.Sprintln(err)
		ctx.SetStatus(500)
		ctx.Write(content)
	}

}

func (ctx *Context) parseFiles(filenames ...string) (*template.Template, error) {
	var t *template.Template = nil

	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}

	for _, filename := range filenames {
		var tmpl *template.Template
		if t == nil {
			t = template.New(filename)
		}
		if filename == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(filename)
		}
		_, err := tmpl.Parse(templates[filename])
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
