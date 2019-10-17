package h2

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
)

func ExamplePost() {
	ts := httptest.NewServer(http.HandlerFunc(a))
	u, err := url.Parse(ts.URL) // ParseRequestURI
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	post1(u)
	ts.Close()

	ts = httptest.NewServer(http.HandlerFunc(b))
	u, err = url.Parse(ts.URL)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	post2(u)
	ts.Close()

	// Output:
	// map[test:[a b]]
	// map[status:ok]
	// "test field value"
	// map[status:ok]
}

func post1(u *url.URL) {
	data := url.Values{}
	data.Set("test", "a")
	data.Add("test", "b")

	r, _ := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	printj(resp.Body)
}

func post2(u *url.URL) {
	r, w := io.Pipe()
	defer r.Close()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		stream, err := m.CreateFormFile("fieldname1", "stream")
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		if _, err = io.Copy(stream, bytes.NewBufferString("test stream data")); err != nil {
			fmt.Fprint(os.Stderr, err)
			return
		}
		f2, err := m.CreateFormField("fieldname2")
		if err != nil {
			log.Fatalln(err)
		}
		_, err = f2.Write([]byte("test field value"))
		if err != nil {
			log.Fatalln(err)
		}
	}()
	resp, err := http.Post(u.String(), m.FormDataContentType(), r)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	printj(resp.Body)
}

func a(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	err := r.ParseForm()
	if err != nil {
		Errorf(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v\n", r.Form)
	// if v, ok := r.Form["foo"]; ok {}
	// fmt.Printf("%s\n", r.FormValue("test"))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func b(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		Errorf(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parse file field
	p, err := reader.NextPart()
	if err != nil && err != io.EOF {
		Errorf(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if p.FormName() != "fieldname1" {
		Errorf(w, "fieldname1 is expected", http.StatusBadRequest)
		return
	}

	buf := bufio.NewReader(p)
	sniff, _ := buf.Peek(512)
	contentType := http.DetectContentType(sniff)
	if contentType != "text/plain; charset=utf-8" {
		Errorf(w, "file type "+contentType+" not allowed", http.StatusBadRequest)
		return
	}

	f := &bytes.Buffer{}
	var maxSize int64 = 32 << 20
	n, err := io.Copy(f, io.MultiReader(buf, io.LimitReader(p, maxSize-511)))
	if err != nil && err != io.EOF {
		Errorf(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n > maxSize {
		Errorf(w, "file size over limit", http.StatusBadRequest)
		return
	}

	// parse text field
	p, err = reader.NextPart()
	if err != nil {
		Errorf(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if p.FormName() != "fieldname2" {
		Errorf(w, "fieldname2 is expected", http.StatusBadRequest)
		return
	}
	text := make([]byte, 64)
	_, err = p.Read(text)
	if err != nil && err != io.EOF {
		Errorf(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("%q\n", string(bytes.Trim(text, "\x00")))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// r.SetBasicAuth("u", "p")
// r.BasicAuth()
// r.Header.Add("Authorization", "Bearer ...")
// r.Header.Get("Authorization")

// http.Redirect(w, r, "/", 302)

// file, fh, err := r.FormFile("fieldname1")
