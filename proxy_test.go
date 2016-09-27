package httx

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestGoHTTP Handler
func TestGoHTTP(t *testing.T) {
	ts1 := httptest.NewServer(http.FileServer(http.Dir("./html")))
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
	ts := httptest.NewServer(http.FileServer(http.Dir("./html")))
	defer ts.Close()
	req := httptest.NewRequest("GET", ts.URL+"/index.html", nil)
	w := httptest.NewRecorder()
	u, _ := url.Parse(ts.URL)
	GoFunc(u, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})(w, req)
	t.Logf("%d - %s", w.Code, w.Body.String())
}
