package httx

import (
	"net/http"
	"os"
)

// RootHTTP Handler
func RootHTTP(root string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := os.Open(root + r.URL.Path)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		r.URL.Path = root + r.URL.Path
		h.ServeHTTP(w, r)
	})
}
