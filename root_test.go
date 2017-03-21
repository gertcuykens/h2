package httx

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestRootHTTP Handler
func TestRootHandler(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tmpfn := filepath.Join(dir, "index.html")
	if err := ioutil.WriteFile(tmpfn, []byte("index.html file's content"), 0666); err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(RootHTTP("./", http.FileServer(http.Dir(dir))))
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
