package appengine

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// ContextHandler for appengine
type ContextHandler struct {
	HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)
}

func (f ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	f.HandlerFunc(c, w, r)
}
