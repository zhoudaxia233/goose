package goose

import "net/http"

// Router represents a collection of path=>handler mappings
type Router struct {
	routers map[string]HandlerFunc
}

func newRouter() *Router {
	return &Router{routers: make(map[string]HandlerFunc)}
}

func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "," + pattern
	r.routers[key] = handler
}

func (r *Router) get(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodGet, pattern, handler)
}

func (r *Router) post(pattern string, handler HandlerFunc) {
	r.addRoute(http.MethodPost, pattern, handler)
}
