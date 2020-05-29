package goose

import "net/http"

// RouterGroup represents a group of routers with the same prefix
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	goose       *Goose
}

func newRouterGroup(goose *Goose) *RouterGroup {
	return &RouterGroup{goose: goose}
}

// Group is used to create a new router group
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	goose := rg.goose
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		goose:  goose,
	}
	goose.groups = append(goose.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	pattern = rg.prefix + pattern
	rg.goose.router.addRoute(method, pattern, handler)
}

// GET is used to handle GET requests
func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.addRoute(http.MethodGet, pattern, handler)
}

// POST is used to handle POST requests
func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.addRoute(http.MethodPost, pattern, handler)
}
