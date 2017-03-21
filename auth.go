package httx

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

// AuthHTTP HandlerFunc
func AuthHTTP(salt string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sum := sha256.Sum256([]byte(r.FormValue("id") + salt))
		if t := r.Header.Get("Authorization"); t != base64.StdEncoding.EncodeToString(sum[:]) {
			http.Error(w, t, http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

// Auth := req.Header.Set("Authorization", auth("111-111-111", "salt"))
func Auth(id string, salt string) string {
	sum := sha256.Sum256([]byte(id + salt))
	return base64.StdEncoding.EncodeToString(sum[:])
}

// example
func secureGET(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fn(w, r)
			return
		}
		if r.Header.Get("Authorization") != "secret" { // r.FormValue("p")
			http.Error(w, "Unauthorized!", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}
