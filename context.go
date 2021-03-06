package goose

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context represents the context of the current HTTP request
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	// request-related

	Path   string
	Method string

	// response-related

	// HTTP status code
	Code int

	// misc

	Params   map[string]string
	handlers []HandlerFunc
	// used in the middleware component
	index int
	goose *Goose
}

func newContext(goose *Goose) *Context {
	return &Context{goose: goose}
}

func (ctx *Context) resetContext(w http.ResponseWriter, r *http.Request) {
	*ctx = Context{
		ResponseWriter: w,
		Request:        r,
		Path:           r.URL.Path,
		Method:         r.Method,
		index:          -1,
		goose:          ctx.goose,
	}
}

// Response part

// Header sets the header information for the response
func (ctx *Context) Header(key, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// StatusCode sets the status code for the response
func (ctx *Context) StatusCode(code int) {
	ctx.Code = code
	ctx.ResponseWriter.WriteHeader(code)
}

// String writes string to the response
func (ctx *Context) String(format string, a ...interface{}) {
	ctx.StringC(http.StatusOK, format, a...)
}

// HTML writes "data" to template "tplName" in the response
func (ctx *Context) HTML(tplName string, data interface{}) {
	ctx.HTMLC(http.StatusOK, tplName, data)
}

// JSON writes json to the response
func (ctx *Context) JSON(obj interface{}) {
	ctx.JSONC(http.StatusOK, obj)
}

// StringC writes string to the response with status Code
func (ctx *Context) StringC(statusCode int, format string, a ...interface{}) {
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	ctx.StatusCode(statusCode)
	if _, err := ctx.ResponseWriter.Write([]byte(fmt.Sprintf(format, a...))); err != nil {
		ctx.Abort(http.StatusInternalServerError, err.Error())
	}
}

// HTMLC writes html to the response with status Code
func (ctx *Context) HTMLC(statusCode int, tplName string, data interface{}) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.StatusCode(statusCode)
	err := ctx.goose.template.templates.ExecuteTemplate(ctx.ResponseWriter, tplName, data)
	if err != nil {
		ctx.Abort(http.StatusInternalServerError, err.Error())
	}
}

// JSONC writes json to the response with status Code
func (ctx *Context) JSONC(statusCode int, obj interface{}) {
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.StatusCode(statusCode)
	encoder := json.NewEncoder(ctx.ResponseWriter)
	if err := encoder.Encode(obj); err != nil {
		ctx.Abort(http.StatusInternalServerError, err.Error())
	}
}

// Request part

// Query returns the value of the given param in the request URL
func (ctx *Context) Query(param string) string {
	return ctx.Request.URL.Query().Get(param)
}

// misc part

// Param returns the value associated with wildcard param in the routing pattern
func (ctx *Context) Param(param string) string {
	value, exists := ctx.Params[param]
	if !exists {
		panic(fmt.Sprintf("Wildcard parameter %s doesn't exist.", param))
	}
	return value
}

// Next gives control to the next handler in Context.handlers
func (ctx *Context) Next() {
	ctx.index++
	numOfHandlers := len(ctx.handlers)
	for ; ctx.index < numOfHandlers; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// Abort stops the execution of current context and write to
// response with statusCode and message
func (ctx *Context) Abort(statusCode int, message string) {
	ctx.index = len(ctx.handlers)
	ctx.JSONC(statusCode, X{"message": message})
}
