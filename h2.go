package h2

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run server
func Run(srv *http.Server) {
	go func(srv *http.Server) {
		var err error
		if srv.TLSConfig != nil {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil {
			fmt.Fprint(os.Stderr, err, "\n")
		}
	}(srv)

	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
	}
}
