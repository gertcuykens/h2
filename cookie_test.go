package h2

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Secure HttpOnly SameSite=Lax|Strict
// Domain, Expires, Max-Age and Path
// Strict-Transport-Security: max-age=3600
// Expect-CT: max-age=3600, enforce, report-uri="https://ct.example.com/report"

func ExampleCookie() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("Flavor"); err != nil {
			http.SetCookie(w, &http.Cookie{Name: "Flavor", Value: "Chocolate Chip"})
		} else {
			cookie.Value = "Oatmeal Raisin"
			http.SetCookie(w, cookie)
		}
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: time.Duration(3 * time.Second),
	}

	if _, err = client.Get(u.String()); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("After 1st request:")
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	if _, err = client.Get(u.String()); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("After 2nd request:")
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	// Output:
	// After 1st request:
	//   Flavor: Chocolate Chip
	// After 2nd request:
	//   Flavor: Oatmeal Raisin
}
