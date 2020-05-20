package goose

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc is the type of request handlers used by goose
type HandlerFunc func(*Context)

// Goose is the core of everything
type Goose struct {
	context *Context
	router  *Router
}

// New is the constructor of goose.Goose
func New() *Goose {
	return &Goose{
		context: newContext(),
		router:  newRouter(),
	}
}

// GET is used to handle GET requests
func (g *Goose) GET(pattern string, handler HandlerFunc) {
	g.router.get(pattern, handler)
}

// POST is used to handle POST requests
func (g *Goose) POST(pattern string, handler HandlerFunc) {
	g.router.post(pattern, handler)
}

func (g *Goose) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.context.resetContext(w, r)
	ctx := g.context
	k := ctx.Method + "," + ctx.Path
	if handler, ok := g.router.routers[k]; ok {
		handler(ctx)
	} else {
		fmt.Fprintf(ctx.ResponseWriter, "404 Not found! - %s\n", ctx.Path)
	}
}

// Run is used to start a http server
func (g *Goose) Run(addr string) error {
	log.Printf("* Running on http://127.0.0.1%s/\n", addr)
	return http.ListenAndServe(addr, g)
}
