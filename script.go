package h2

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// ScriptHTTP Handler
func ScriptHTTP(script string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			h.ServeHTTP(w, r)
			return
		}
		f, err := ioutil.ReadFile("index.html")
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		s := string(f)
		v := strings.Replace(s, "</body>", script+"</body>", 1)
		w.Write([]byte(v))
	})
}
