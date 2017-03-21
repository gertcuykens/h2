package httx

import (
	"net/http"
)

// CorsHandler httx
func CorsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(CorsHandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}))
}

// CorsHandlerFunc httx
func CorsHandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")
		//w.Header().Set("Access-Control-Allow-Credentials", "true")
		switch r.Method {
		case "OPTIONS":
		default:
			fn(w, r)
		}
	}
}
