package meiporo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/lestrrat/go-xslate"
)

type ResponseWriter interface {
	http.ResponseWriter
	Status() int
	Written() bool
	Size() int
	WriteJSON(interface{})
	WriteXML(interface{})
	WriteHTML(string, ...interface{})
	WriteText(string, ...interface{})
	Render(string, xslate.Vars)
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func NewResponseWriter(res http.ResponseWriter) ResponseWriter {
	return &responseWriter{res, 0, 0}
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) WriteJSON(data interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(data)
}

func (rw *responseWriter) WriteXML(data interface{}) {
	rw.Header().Set("Content-Type", "application/xml; charset=utf-8")
	xml.NewEncoder(rw).Encode(data)
}

func (rw *responseWriter) WriteHTML(textFormat string, data ...interface{}) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(rw, textFormat, data...)
}

func (rw *responseWriter) WriteText(textFormat string, data ...interface{}) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(rw, textFormat, data...)
}

func (rw *responseWriter) Render(path string, data xslate.Vars) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	output, _ := Xslate.Render(path, data)
	fmt.Fprintf(rw, output)
}
