package httx

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRootHTTP Handler
func TestRootHandler(t *testing.T) {
	ts := httptest.NewServer(RootHTTP("./html", http.FileServer(http.Dir("./"))))
	defer ts.Close()
	res, err := http.Get(ts.URL + "/index.html")
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
