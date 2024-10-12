package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipMiddleware struct{}

func NewGzipMiddleware() *GzipMiddleware {
	return &GzipMiddleware{}
}

func (g *GzipMiddleware) GzipMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {

			// create a gziped response

			wrw := NewWrappedResponseWriter(rw)

			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, r)
			defer wrw.Flush()

			return

		}

		// If the client does not support gzip then don't gzip

		next.ServeHTTP(rw, r)

	})
}

type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {

	gw := gzip.NewWriter(rw)

	return &WrappedResponseWriter{rw: rw, gw: gw}
}

func (wrw *WrappedResponseWriter) Header() http.Header {
	return wrw.rw.Header()
}

func (wrw *WrappedResponseWriter) Write(p []byte) (int, error) {
	return wrw.gw.Write(p)
}

func (wrw *WrappedResponseWriter) WriteHeader(statusCode int) {
	wrw.rw.WriteHeader(statusCode)
}

func (wrw *WrappedResponseWriter) Flush() {
	wrw.gw.Flush()
	wrw.gw.Close()
}
