package ela

import (
	// "github.com/gogather/com/log"
	"html/template"
	"net/http"
	"path/filepath"
	// "os"
)

// RequestContext
type RequestContext struct {
	w         http.ResponseWriter
	r         *http.Request
	Data      map[string]interface{}
	status    int
	headerMap map[string]string
}

func (this *RequestContext) GetResponseWriter() http.ResponseWriter {
	return this.w
}

func (this *RequestContext) GetRequest() *http.Request {
	return this.r
}

func (this *RequestContext) SetStatus(status int) {
	this.status = status
}

func (this *RequestContext) SetHeader(key, value string) {
	if this.headerMap == nil {
		this.headerMap = make(map[string]string)
	}
	this.headerMap[key] = value
}

func (this *RequestContext) Write(content string) (int, error) {
	this.writeHeader()
	return this.w.Write([]byte(content))
}

func (this *RequestContext) ServeTemplate(templateFile string) {
	this.SetHeader("Content-Type", "text/html")
	this.writeHeader()

	templateFile = filepath.Join("view", templateFile)
	t, _ := template.ParseFiles(templateFile)
	t.Execute(this.w, this.Data)
}

func (this *RequestContext) writeHeader() {
	for k, v := range this.headerMap {
		this.r.Header.Set(k, v)
	}
	this.w.WriteHeader(this.status)
}
