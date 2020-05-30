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

func (r *Router) getRoute(method string, pattern string) (searchResultPtr *node, params map[string]string) {
	searchResultPtr, params = r.routers[method].search(pattern)
	return
}

func (r *Router) handleRequest(ctx *Context) {
	searchResultPtr, params := r.getRoute(ctx.Method, ctx.Path)

	handler := searchResultPtr.handler
	ctx.Params = params

	if handler == nil {
		handler = func(ctx *Context) {
			ctx.setString(404, "404 Not Found! - %s", ctx.Path)
		}
	}
	ctx.handlers = append(ctx.handlers, handler)
	ctx.Next()
}
