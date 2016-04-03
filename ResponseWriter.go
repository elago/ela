package ela

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (this *responseWriter) Status() int {
	return this.status
}

func (this *responseWriter) Header() http.Header {
	return this.ResponseWriter.Header()
}

func (this *responseWriter) Write(data []byte) (int, error) {
	return this.ResponseWriter.Write(data)
}

func (this *responseWriter) WriteHeader(status int) {
	this.ResponseWriter.WriteHeader(status)
	this.status = status
}
