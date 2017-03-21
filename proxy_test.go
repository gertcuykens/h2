package httx

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

// TestGoHTTP Handler
func TestGoHTTP(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tmpfn := filepath.Join(dir, "index.go")
	if err := ioutil.WriteFile(tmpfn, []byte("index.go file's content"), 0666); err != nil {
		t.Fatal(err)
	}

	ts1 := httptest.NewServer(http.FileServer(http.Dir(dir)))
	defer ts1.Close()
	u, _ := url.Parse(ts1.URL)

	ts2 := httptest.NewServer(GoHTTP(u, http.FileServer(http.Dir("./"))))
	defer ts2.Close()

	res, err := http.Get(ts2.URL + "/index.go")
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%d - %s", res.StatusCode, body)
}

// TestGoFunc HandlerFunc
func TestGoFunc(t *testing.T) {
	req := httptest.NewRequest("GET", "/index.html", nil)
	w := httptest.NewRecorder()
	u, _ := url.Parse("/index.html")
	GoFunc(u, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})(w, req)
	t.Logf("%d - %s", w.Code, w.Body.String())
}
