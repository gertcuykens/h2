package httx

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRootHTTP Handler
func TestRootHandler(t *testing.T) {
	ts := httptest.NewServer(RootHTTP(".", http.FileServer(http.Dir("./httx_test"))))
	defer ts.Close()
	res, err := http.Get(ts.URL)
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
