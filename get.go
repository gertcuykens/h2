package h2

import (
	"io/ioutil"
	"net/http"
)

func Get(r *http.Request, f func([]byte)) error {
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	f(b)
	if c, ok := res.Header["Set-Cookie"]; ok {
		r.Header = http.Header{"Cookie": c}
	}
	return nil
}
