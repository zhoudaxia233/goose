package goose

import (
	"net/http"
	"testing"
)

func TestAsteroidRouting(t *testing.T) {
	signature := ""
	g := New()

	g.GET("/assets/*files", func(ctx *Context) {
		signature += ctx.Param("*files")
	})

	// test /assets/hello
	signature = ""
	w := sendRequest(g, http.MethodGet, "/assets/hello")
	if w.Code != http.StatusOK {
		t.Errorf("Got status code %d instead of 200.", w.Code)
	}
	want := "hello"
	got := signature
	if got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}

	// test /assets/this/is/path/to/static/files/like/hello.js
	signature = ""
	w = sendRequest(g, http.MethodGet, "/assets/this/is/path/to/static/files/like/hello.js")
	if w.Code != http.StatusOK {
		t.Errorf("Got status code %d instead of 200.", w.Code)
	}
	want = "this/is/path/to/static/files/like/hello.js"
	got = signature
	if got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}
}
