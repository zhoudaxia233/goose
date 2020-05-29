package goose

import (
	"log"
	"net/http"
)

// HandlerFunc is the type of request handlers used by goose
type HandlerFunc func(*Context)

// Goose is a top-level framework instance
type Goose struct {
	*RouterGroup
	groups  []*RouterGroup
	context *Context
	router  *Router
}

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

	searchResultPtr, params := g.router.routers[ctx.Method].search(ctx.Path)
	ctx.Params = params
	handler := searchResultPtr.handler

	if handler != nil {
		handler(ctx)
	} else {
		ctx.setString(404, "404 Not Found! - %s", ctx.Path)
	}
}

// Run is used to start a http server
func (g *Goose) Run(addr string) error {
	log.Printf("* Running on http://127.0.0.1%s/\n", addr)
	return http.ListenAndServe(addr, g)
}
