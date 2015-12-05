package ela

import (
	"html/template"
	"net/http"
	// "os"
)

// RequestContext
type RequestContext struct {
	w    http.ResponseWriter
	r    *http.Request
	Data map[string]interface{}
}

func (this *RequestContext) Write(content string) (int, error) {
	return this.w.Write([]byte(content))
}

func (this *RequestContext) ServeTemplate(templateFile string) {
	t := template.New("fieldname example")
	t, _ = t.Parse("hello {{.name}}! {{if eq .id -1}}below zero{{end}}")

	t.Execute(this.w, this.Data)
}
