package ela

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	SetStatus(int)
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (r *responseWriter) SetStatus(status int) {
	r.status = status
}

func (r *responseWriter) Status() int {
	return r.status
}

func (r *responseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *responseWriter) Write(data []byte) (int, error) {
	return r.ResponseWriter.Write(data)
}

func (r *responseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}
