package goose

import (
	"net/http"
)

// Router represents a collection of path=>handler mappings
type Router struct {
	routers map[string]*node
}

func newRouter() *Router {
	return &Router{routers: make(map[string]*node)}
}

func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	if _, exists := r.routers[method]; !exists {
		r.routers[method] = &node{}
	}
	r.routers[method].insert(pattern, handler)
}

func (r *Router) get(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodGet, pattern, handler)
}

func (r *Router) post(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodPost, pattern, handler)
}

func (r *Router) handle(ctx *Context) {
	handler := r.routers[ctx.Method].search(ctx.Path).handler

	if handler != nil {
		handler(ctx)
	} else {
		ctx.setString(404, "404 Not Found! - ", ctx.Path)
	}
}
