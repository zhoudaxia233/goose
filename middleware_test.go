package goose

import (
	"net/http"
	"testing"
)

func TestCustomMiddleware(t *testing.T) {
	signature := ""
	g := New()
	g.Use(func(ctx *Context) {
		signature += "a"
		ctx.Next()
		signature += "b"
	})

	// I put it here to show that the relative position of method Use and GET doesn't matter
	g.GET("/", func(ctx *Context) {
		signature += "/"
	})

	g.Use(func(ctx *Context) {
		signature += "c"
		ctx.Next()
		signature += "d"
	})

	g.Use(func(ctx *Context) {
		signature += "e"
	})

	v1 := g.Group("v1")
	v1.Use(func(ctx *Context) {
		signature += "f"
		ctx.Next()
		signature += "g"
	})

	v1.GET("/", func(ctx *Context) {
		signature += "@"
	})

	v1.GET("/hello", func(ctx *Context) {
		signature += "#"
	})

	// test /
	signature = ""
	w := sendRequest(g, http.MethodGet, "/")
	if w.Code != http.StatusOK {
		t.Errorf("Got status code %d instead of 200.", w.Code)
	}
	want := "ace/db"
	got := signature
	if got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}

	// test /v1/
	signature = ""
	w = sendRequest(g, http.MethodGet, "/v1/")
	if w.Code != http.StatusOK {
		t.Errorf("Got status code %d instead of 200.", w.Code)
	}
	want = "acef@gdb"
	got = signature
	if got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}

	// test /v1/hello
	signature = ""
	w = sendRequest(g, http.MethodGet, "/v1/hello")
	if w.Code != http.StatusOK {
		t.Errorf("Got status code %d instead of 200.", w.Code)
	}
	want = "acef#gdb"
	got = signature
	if got != want {
		t.Errorf("Want: %s , Got: %s", want, got)
	}
}
