package h2

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// type roundTripper func(*http.Request) (*http.Response, error)

// func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func proxy(u *url.URL) http.Handler {
	return &httputil.ReverseProxy{
		// Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
		// 	req.URL.Scheme = u.Scheme
		// 	req.URL.Host = u.Host
		// 	req.URL.Path = "/authorization/" + trim(req.URL.Path)
		// 	req.Header.Set("Host", u.Host)
		// 	req.Host = u.Host
		// 	return http.DefaultTransport.RoundTrip(req)
		// }),
		Director: func(req *http.Request) {
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			req.URL.Path = "/authorization/" + trim(req.URL.Path)
			req.Header.Set("Host", u.Host)
			req.Host = u.Host
		},
	}
}

func trim(p string) string {
	i := bytes.IndexByte([]byte(p), '/')
	if i > -1 {
		return trim(p[i+1:])
	}
	return p
}

func ProxyMux(x *http.ServeMux, aud []string) {
	for _, a := range aud {
		u, err := url.Parse(a)
		if err != nil {
			panic(err)
		}
		x.Handle("/"+u.Host+"/", proxy(u))
	}
}
