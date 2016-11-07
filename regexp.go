package httx

import (
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// RegexpHTTP routes
type RegexpHTTP struct {
	routes []*route
}

// Handle RegexpHTTP
func (h *RegexpHTTP) Handle(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

// HandleFunc RegexpHTTP
func (h *RegexpHTTP) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

// ServeHTTP RegexpHTTP
func (h *RegexpHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		// strings.Index(r.URL.Path, route.pattern)
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
