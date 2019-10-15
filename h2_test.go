package h2

import (
	"flag"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func ExampleRun() {
	var host string
	flag.StringVar(&host, "host", "", "domain name")
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(""))) // http.StripPrefix(p, h)
	mux.HandleFunc("/tls", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "TLS: %+v", r.TLS)
	})

	if len(host) > 0 {
		m := autocert.Manager{
			Cache:      autocert.DirCache("tls"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(host),
		}
		srv := &http.Server{
			Addr:    ":https",
			Handler: mux,
			// Handler: m.HTTPHandler(nil),
			// TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
			TLSConfig: m.TLSConfig(),
		}
		fmt.Println("https://" + host)
		go Run(srv)
	} else {
		srv := &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		fmt.Println("http://localhost:8080")
		go Run(srv)
	}
	time.Sleep(1 * time.Second)
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	// Output:
	// http://localhost:8080
}
