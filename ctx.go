package ela

import (
	"fmt"
	"github.com/gogather/com/log"
	"html/template"
	"net/http"
)

// RequestContext
type Context struct {
	w         http.ResponseWriter
	r         *http.Request
	Data      map[string]interface{}
	status    int
	headerMap map[string]string
}

func (this *Context) GetResponseWriter() http.ResponseWriter {
	return this.w
}

func (this *Context) GetRequest() *http.Request {
	return this.r
}

func (this *Context) GetMethod() string {
	return this.r.Method
}

func (this *Context) SetStatus(status int) {
	this.status = status
}

func (this *Context) GetStatus() int {
	return this.status
}

func (this *Context) SetHeader(key, value string) {
	if this.headerMap == nil {
		this.headerMap = make(map[string]string)
	}
	this.headerMap[key] = value
}

func (this *Context) GetCookie(key string) (*http.Cookie, error) {
	return this.r.Cookie(key)
}

func (this *Context) SetCookie(cookie *http.Cookie) {
	log.Redln(cookie)
	http.SetCookie(this.w, cookie)
}

func (this *Context) Write(content string) (int, error) {
	this.writeHeader()
	return this.w.Write([]byte(content))
}

func (this *Context) ServeTemplate(templateFile string) {
	this.SetHeader("Content-Type", "text/html")
	this.writeHeader()

	ReloadTemplate() // TODO: if in debug mode, reload templates

	t, err := this.parseFiles(templatesName...)

	if err != nil {
		log.Redln(err)
		// this.w.WriteHeader(500)
	} else {
		err = t.ExecuteTemplate(this.w, templateFile, this.Data)
		// log.Pinkln(err)
	}

}

func (this *Context) parseFiles(filenames ...string) (*template.Template, error) {
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

func (this *Context) writeHeader() {
	for k, v := range this.headerMap {
		this.r.Header.Set(k, v)
	}
	this.w.WriteHeader(this.status)
}
