package ela

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	SetStatus(int)
	HasFlushed() bool
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0, false}
}

type responseWriter struct {
	http.ResponseWriter
	status  int
	flushed bool
}

func (r *responseWriter) SetStatus(status int) {
	r.status = status
}

func (r *responseWriter) Status() int {
	return r.status
}

func (r *responseWriter) HasFlushed() bool {
	return r.flushed == true
}

func (r *responseWriter) Header() http.Header {
	r.flushed = true
	return r.ResponseWriter.Header()
}

func (r *responseWriter) Write(data []byte) (int, error) {
	return r.ResponseWriter.Write(data)
}

func (r *responseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}
