package h2

import (
	"crypto/subtle"
	"net/http"
)

func BasicAuth(fn http.HandlerFunc, usr, pwd string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(usr+pwd), []byte(u+p)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+r.Referer()+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`"unauthorised"`))
			return
		}
		fn(w, r)
	}
}
