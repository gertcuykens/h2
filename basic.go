package h2

import (
	"crypto/subtle"
	"net/http"
)

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		usr, pwd, ok := r.BasicAuth()

		if !ok ||
			subtle.ConstantTimeCompare([]byte(usr), []byte(username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pwd), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			// jsonResponse(w, "Unauthorised", http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}
