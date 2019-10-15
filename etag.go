package h2

import (
	"net/http"
	"strings"
)

// EtagHTTP Handler
func EtagHTTP(etag string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("If-None-Match"), etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Cache-Control", "max-age=120")
		w.Header().Set("ETag", etag)
		h.ServeHTTP(w, r)
	})
}
