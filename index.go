package main

import (
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

// func init() {
// 	http.HandleFunc("/", gzipFunc(index))
// }

func index(w http.ResponseWriter, r *http.Request) {
	var push = [][]string{
		[]string{"favicon.ico", "image/x-icon"},
		[]string{"manifest.json", "application/javascript"},
		[]string{"index.html", "text/html; charset=utf-8"},
		[]string{"index.css", "text/css; charset=utf-8"},
	}
	for k := range push {
		f, err := ioutil.ReadFile(push[k][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sha := sha1.Sum(f)
		etag := base64.URLEncoding.EncodeToString(sha[:])
		w.Header().Set("ETag", etag)
		w.Header().Set("Content-Type", push[k][1])
		if p, ok := w.(http.Pusher); ok {
			p.Push(push[k][0], nil)
		}
	}
}
