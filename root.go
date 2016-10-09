package httx

import (
	"net/http"
	"os"
)

// RootHTTP Handler
func RootHTTP(root string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			h.ServeHTTP(w, r)
			return
		}
		if _, err := os.Stat(root + r.URL.Path); os.IsNotExist(err) {
			h.ServeHTTP(w, r)
			return
		}
		h2 := http.FileServer(http.Dir(root))
		h2.ServeHTTP(w, r)
	})
}
