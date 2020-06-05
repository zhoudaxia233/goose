package goose

import (
	tpl "html/template"
	"log"
	"net/http"
	"strings"
)

type (
	// HandlerFunc is the type of request handlers used by goose
	HandlerFunc func(*Context)

	// X is a shortcut for map[string]interface{}
	X map[string]interface{}

	// Template acts as an interface of Go's html templating system used in goose
	Template struct {
		funcMap   X
		templates *tpl.Template
	}

	// Goose is a top-level framework instance
	Goose struct {
		*RouterGroup
		groups   []*RouterGroup
		context  *Context
		router   *Router
		template *Template
	}
)

// New is the constructor of goose.Goose
func New() *Goose {
	goose := &Goose{
		router:   newRouter(),
		template: &Template{funcMap: make(X)},
	}
	goose.context = newContext(goose)
	goose.RouterGroup = newRouterGroup(goose)
	goose.groups = []*RouterGroup{goose.RouterGroup}
	return goose
}

// FuncMap takes a map of type X and use it as the funcMap of goose
func (g *Goose) FuncMap(funcMap X) {
	g.template.funcMap = funcMap
}

// Set sets the funcMap entries which defines the mapping from names to functions
func (g *Goose) Set(key string, value interface{}) {
	g.template.funcMap[key] = value
}

// LoadHTMLGlob loads a glob of HTML templates in one go
func (g *Goose) LoadHTMLGlob(pattern string) {
	g.template.templates = tpl.Must(tpl.New("").Funcs(tpl.FuncMap(g.template.funcMap)).ParseGlob(pattern))
}

func (g *Goose) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := g.context
	ctx.resetContext(w, r)

	g.applyMiddlewares()
	g.router.handleRequest(ctx)
}

// Run is used to start a http server
func (g *Goose) Run(addr string) error {
	log.Printf("* Running on http://127.0.0.1%s/\n", addr)
	return http.ListenAndServe(addr, g)
}

func (g *Goose) applyMiddlewares() {
	ctx := g.context
	middlewares := []HandlerFunc{}
	for _, group := range g.groups {
		if strings.HasPrefix(ctx.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	ctx.handlers = middlewares
}
