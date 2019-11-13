package h2

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func ExampleRun() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(""))) // http.StripPrefix(p, h)
	mux.HandleFunc("/tls", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TLS: %+v", r.TLS)
	})

	u, err := url.Parse(os.Getenv("ORIGIN"))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    ":" + u.Port(),
		Handler: mux,
	}

	if u.Scheme == "https" {
		m := autocert.Manager{
			Cache:      autocert.DirCache("tls"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(u.Host),
		}
		srv.TLSConfig = m.TLSConfig()
		// Handler: m.HTTPHandler(nil),
		// TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}
	fmt.Println(os.Getenv("ORIGIN"))
	go Run(srv)
	time.Sleep(1 * time.Second)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	// Output:
	// http://localhost:8080
}
