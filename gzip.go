package h2

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHTTP Handler
func GzipHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(GzipFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}))
}

// GzipFunc HandlerFunc
func GzipFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		g := gzip.NewWriter(w)
		defer g.Close()
		z := gzipResponseWriter{Writer: g, ResponseWriter: w}
		fn(z, r)
	}
}
