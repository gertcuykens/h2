package h2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// our service
type CakeBaker struct{}

func (CakeBaker) BakeIt(n int, msg *string) error {
	*msg = fmt.Sprintf("your cake has been baked (%d)", n)
	return nil
}

type Args struct {
	A, B int
}

func ExampleCodec() {
	srv := rpc.NewServer()
	srv.Register(new(CakeBaker))

	mux := http.NewServeMux()
	mux.HandleFunc("/bake", func(w http.ResponseWriter, r *http.Request) {
		jsonCodec := jsonrpc.NewServerCodec(Codec{i: r.Body, o: w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err := srv.ServeRequest(jsonCodec)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/bake", "application/json", bytes.NewBufferString(
		`{"id":1,"method":"CakeBaker.BakeIt","params":[10]}`,
	))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", string(b))
	// Output:
	// {"id":1,"result":"your cake has been baked (10)","error":null}
}

func ExampleClient() {
	srv := rpc.NewServer()
	srv.Register(new(CakeBaker))

	go func(srv *rpc.Server) {
		l, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			go func(conn net.Conn) {
				srv.ServeConn(conn)
				defer conn.Close()
			}(conn)
		}
	}(srv)

	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var reply string
	err = client.Call("CakeBaker.BakeIt", 10, &reply)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer client.Close()
	fmt.Printf("%s\n", reply)

	// Output:
	// your cake has been baked (10)
}
