package goose

import (
	"log"
	"net/http"
	"os"
	"path"
)

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
		prefix: path.Join("/", rg.prefix, prefix),
		goose:  goose,
	}
	goose.groups = append(goose.groups, newGroup)
	return newGroup
}

// Use is used to add middlewares to a router group
func (rg *RouterGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
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

// StaticFile registers a router for serving a single static file, such as favicon.ico
func (rg *RouterGroup) StaticFile(pattern, nativePath string) {
	checkExistence(nativePath)
	handler := func(ctx *Context) {
		http.ServeFile(ctx.ResponseWriter, ctx.Request, nativePath)
	}
	rg.GET(pattern, handler)
}

// Static registers a router for serving static files
func (rg *RouterGroup) Static(pattern, nativePath string) {
	checkExistence(nativePath)
	handler := rg.makeStaticHandler(pattern, http.Dir(nativePath))
	urlPattern := path.Join(pattern, "/*files")
	rg.GET(urlPattern, handler)
}

func (rg *RouterGroup) makeStaticHandler(pattern string, fileSystem http.FileSystem) HandlerFunc {
	pattern = path.Join(rg.prefix, pattern)
	fileServer := http.StripPrefix(pattern, http.FileServer(fileSystem))

	return func(ctx *Context) {
		file := ctx.Param("*files")
		// check if file exists and/or if we have the permission to access it
		/* the reason we add this code block is: fileServer is actually handled by
		   the official http library. Therefore, when we access a file that either
		   doesn't exist or permission-denied, our custom errorHandler doesn't work.
		   http's error handler will take it over. And this is normally an unexpected
		   behavior to the user of goose framework.
		*/
		f, err := fileSystem.Open(file)
		if err != nil {
			ctx.StatusCode(http.StatusNotFound)
			return
		}
		f.Close()

		fileServer.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
}

func checkExistence(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Warning: %s doesn't exist.\n", path)
	}
}
