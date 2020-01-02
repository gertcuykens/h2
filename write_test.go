package h2

import (
	"fmt"
	"net/http/httptest"
)

func ExampleJsonResponse() {
	w := httptest.NewRecorder()
	JsonResponse(w, "test", 200)
	fmt.Printf("%d - %s\n", w.Code, w.Body.String())
	w = httptest.NewRecorder()
	JsonResponse(w, "test", 501)
	fmt.Printf("%d - %s\n", w.Code, w.Body.String())
	// Output:
	// 200 - "test"
	// 501 - "test"
}
