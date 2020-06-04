package goose

import (
	"log"
	"net/http"
	"strings"
)

type (
	// HandlerFunc is the type of request handlers used by goose
	HandlerFunc func(*Context)

	// X is a shortcut for map[string]interface{}
	X map[string]interface{}

	// Goose is a top-level framework instance
	Goose struct {
		*RouterGroup
		groups  []*RouterGroup
		context *Context
		router  *Router
	}
)

// New is the constructor of goose.Goose
func New() *Goose {
	goose := &Goose{
		context: newContext(),
		router:  newRouter(),
	}
	goose.RouterGroup = newRouterGroup(goose)
	goose.groups = []*RouterGroup{goose.RouterGroup}
	return goose
}

func (g *Goose) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := g.context
	ctx.resetContext(w, r)

	g.applyMiddlewares()
	g.router.handleRequest(ctx)
}

// Run is used to start a http server
func (g *Goose) Run(addr string) error {
	log.Printf("* Running on http://127.0.0.1%s/\n", addr)
	return http.ListenAndServe(addr, g)
}

func (g *Goose) applyMiddlewares() {
	ctx := g.context
	middlewares := []HandlerFunc{}
	for _, group := range g.groups {
		if strings.HasPrefix(ctx.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	ctx.handlers = middlewares
}
