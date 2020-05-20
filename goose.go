package goose

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc is the type of request handlers used by goose
type HandlerFunc func(*Context)

// Mux is an HTTP request multiplexer
type Mux struct {
	router map[string]HandlerFunc
}

// New is the constructor of goose.Mux
func New() *Mux {
	return &Mux{router: make(map[string]HandlerFunc)}
}

func (m *Mux) addRoute(method string, pattern string, handler HandlerFunc) {
	k := method + pattern
	m.router[k] = handler
}

// GET is used to handle GET requests
func (m *Mux) GET(pattern string, handler HandlerFunc) {
	m.addRoute("GET", pattern, handler)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(w, r)
	k := ctx.Method + ctx.Path
	if handler, ok := m.router[k]; ok {
		handler(ctx)
	} else {
		fmt.Fprintf(ctx.ResponseWriter, "404 Not found! - %s\n", ctx.Path)
	}
}

// Run is used to start a http server
func (m *Mux) Run(addr string) error {
	log.Printf("* Running on http://127.0.0.1%s/\n", addr)
	return http.ListenAndServe(addr, m)
}
