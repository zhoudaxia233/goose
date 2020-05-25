package goose

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H is a shortcut for map[string]interface{}, it aims to facilitate the construction of JSON objects
type H map[string]interface{}

// Context represents the context of the current HTTP request
type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	// request-related
	Path   string
	Method string

	// response-related
	StatusCode int
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

// Response part

// SetHeader sets the header information for the response
func (ctx *Context) SetHeader(key, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// SetStatusCode sets the status code for the response
func (ctx *Context) SetStatusCode(statusCode int) {
	ctx.StatusCode = statusCode
	ctx.ResponseWriter.WriteHeader(statusCode)
}

// String writes string to the response
func (ctx *Context) String(format string, a ...interface{}) {
	ctx.setString(http.StatusOK, format, a...)
}

// HTML writes html to the response
func (ctx *Context) HTML(html string) {
	ctx.setHTML(http.StatusOK, html)
}

// JSON writes json to the response
func (ctx *Context) JSON(obj interface{}) {
	ctx.setJSON(http.StatusOK, obj)
}

func (ctx *Context) setString(statusCode int, format string, a ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain; charset=utf-8")
	ctx.SetStatusCode(statusCode)
	if _, err := ctx.ResponseWriter.Write([]byte(fmt.Sprintf(format, a...))); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func (ctx *Context) setHTML(statusCode int, html string) {
	ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
	ctx.SetStatusCode(statusCode)
	if _, err := ctx.ResponseWriter.Write([]byte(html)); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func (ctx *Context) setJSON(statusCode int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)
	encoder := json.NewEncoder(ctx.ResponseWriter)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
	}
}

// Request part

// Query returns the value of the given param in the request URL
func (ctx *Context) Query(param string) string {
	return ctx.Request.URL.Query().Get(param)
}
