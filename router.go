package goose

// Router represents a collection of "HTTP verbs=>routing tree" mappings
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
