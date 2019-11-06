package h2

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

func printj(r io.ReadCloser) {
	if r == nil {
		fmt.Fprint(os.Stderr, "reader is nil")
		return
	}
	defer r.Close()
	var j map[string]interface{}
	err := json.NewDecoder(r).Decode(&j)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Printf("%v\n", j)
}

func printb(r io.ReadCloser) {
	if r == nil {
		fmt.Fprint(os.Stderr, "reader is nil")
		return
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	fmt.Printf("%s\n", string(b))
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Errorf(w http.ResponseWriter, error string, code int) {
	msg := fmt.Sprintf(`{"code":%d,"status":"%s","error":%q}`, code, http.StatusText(code), error)
	http.Error(w, msg, code)
}

func d1(w http.ResponseWriter, r *http.Request) {
	printb(http.MaxBytesReader(w, r.Body, 1024))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", `{"code":200,"status":"OK"}`)
}

func d2(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	fmt.Printf("%+v\n", r.Form)
	fmt.Printf("%+v\n", r.MultipartForm)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func d3(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%q", dump)
}

func jsonResponse(w http.ResponseWriter, d interface{}, c int) {
	j, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		http.Error(w, "Error marshal JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	w.Write(j)
}

func jsonStreamResponse(w http.ResponseWriter, d io.Reader, c int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err := enc.Encode(d)
	if err != nil {
		http.Error(w, "Error encode JSON response", http.StatusInternalServerError)
	}
}