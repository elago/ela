package ela

import (
	// "github.com/gogather/com/log"
	"html/template"
	"net/http"
	// "path/filepath"
	// "os"
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

func (this *Context) SetHeader(key, value string) {
	if this.headerMap == nil {
		this.headerMap = make(map[string]string)
	}
	this.headerMap[key] = value
}

func (this *Context) Write(content string) (int, error) {
	this.writeHeader()
	return this.w.Write([]byte(content))
}

func (this *Context) ServeTemplate(templateFile string) {
	this.SetHeader("Content-Type", "text/html")
	this.writeHeader()

	t, err := template.ParseFiles(templatesName...)
	if err != nil {
		this.w.WriteHeader(500)
	} else {
		t.ExecuteTemplate(this.w, templateFile, this.Data)
	}

}

func (this *Context) writeHeader() {
	for k, v := range this.headerMap {
		this.r.Header.Set(k, v)
	}
	this.w.WriteHeader(this.status)
}
