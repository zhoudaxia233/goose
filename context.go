package goose

import "net/http"

// Context represents the context of the current HTTP request
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	// request-related
	Path   string
	Method string
}

func newContext() *Context {
	return &Context{}
}

func (ctx *Context) resetContext(w http.ResponseWriter, r *http.Request) {
	*ctx = Context{
		ResponseWriter: w,
		Request:        r,
		Path:           r.URL.Path,
		Method:         r.Method,
	}
}

// Query returns the value of the given param in the request URL
func (ctx *Context) Query(param string) string {
	return ctx.Request.URL.Query().Get(param)
}
