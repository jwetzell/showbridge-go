package routing

import (
	"net/http"
	"sync"
)

type httpServerContextKey string

const HttpServerResponseContextKey = httpServerContextKey("responseWriter")

func NewHTTPServerResponseWriter(responseWriter http.ResponseWriter) *HTTPServerResponseWriter {
	return &HTTPServerResponseWriter{
		responseWriter: responseWriter,
		done:           false,
	}
}

type HTTPServerResponseWriter struct {
	responseWriter http.ResponseWriter
	done           bool
	mu             sync.Mutex
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
	hsrw.done = true
	hsrw.responseWriter.WriteHeader(status)
}

func (hsrw *HTTPServerResponseWriter) Write(data []byte) (int, error) {
	hsrw.done = true
	return hsrw.responseWriter.Write(data)
}

func (hsrw *HTTPServerResponseWriter) Done() bool {
	return hsrw.done
}
