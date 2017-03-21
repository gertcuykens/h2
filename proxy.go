package httx

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GoHTTP Handler
func GoHTTP(u *url.URL, h http.Handler) http.Handler {
	return http.HandlerFunc(GoFunc(u,
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}))
}

// GoFunc HandlerFunc
func GoFunc(u *url.URL, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Index(r.URL.Path, ".go") == -1 && strings.Index(r.URL.Path, ".csv") == -1 {
			// fmt.Printf("%s\n", r.URL.Path)
			fn(w, r)
			return
		}
		client := &http.Client{}
		q := "?"
		if r.URL.RawQuery == "" {
			q = ""
		}
		fmt.Printf("%s://%s%s%s%s\n", u.Scheme, u.Host, r.URL.Path, q, r.URL.RawQuery)
		req, err := http.NewRequest(r.Method, fmt.Sprintf("%s://%s%s%s%s", u.Scheme, u.Host, r.URL.Path, q, r.URL.RawQuery), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		for k, v := range response.Header {
			w.Header().Set(k, v[0])
		}
		w.WriteHeader(response.StatusCode)
		io.Copy(w, response.Body)
	}
}
