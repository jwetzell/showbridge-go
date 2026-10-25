package routing

import (
	"io"
	"net/http"
	"sync"
)

type contextKey string

const HttpServerResponseContextKey = contextKey("responseWriter")
const WriterContextKey = contextKey("writer")

type ContextWriter struct {
	writer io.Writer
	done   bool
	mu     sync.RWMutex
}

func NewContextWriter(writer io.Writer) *ContextWriter {
	return &ContextWriter{
		writer: writer,
		done:   false,
	}
}

func (cw *ContextWriter) TryLock() bool {
	return cw.mu.TryLock()
}

func (cw *ContextWriter) Lock() {
	cw.mu.Lock()
}

func (cw *ContextWriter) Unlock() {
	cw.mu.Unlock()
}

func (cw *ContextWriter) Done() bool {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	return cw.done
}

func (cw *ContextWriter) Write(data []byte) (int, error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.done = true
	return cw.writer.Write(data)
}

func NewHTTPServerResponseWriter(responseWriter http.ResponseWriter) *HTTPServerResponseWriter {
	return &HTTPServerResponseWriter{
		responseWriter: responseWriter,
		done:           false,
	}
}

type HTTPServerResponseWriter struct {
	responseWriter http.ResponseWriter
	done           bool
	mu             sync.RWMutex
}

func (hsrw *HTTPServerResponseWriter) TryLock() bool {
	return hsrw.mu.TryLock()
}

func (hsrw *HTTPServerResponseWriter) Lock() {
	hsrw.mu.Lock()
}

func (hsrw *HTTPServerResponseWriter) Unlock() {
	hsrw.mu.Unlock()
}

func (hsrw *HTTPServerResponseWriter) WriteHeader(status int) {
	hsrw.mu.Lock()
	defer hsrw.mu.Unlock()
	hsrw.done = true
	hsrw.responseWriter.WriteHeader(status)
}

func (hsrw *HTTPServerResponseWriter) Write(data []byte) (int, error) {
	hsrw.mu.Lock()
	defer hsrw.mu.Unlock()
	hsrw.done = true
	return hsrw.responseWriter.Write(data)
}

func (hsrw *HTTPServerResponseWriter) Done() bool {
	hsrw.mu.RLock()
	defer hsrw.mu.RUnlock()
	return hsrw.done
}
