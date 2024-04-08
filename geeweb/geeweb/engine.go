package geeweb

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type MethodType string

const (
	GET  MethodType = "GET"
	POST MethodType = "POST"
)

// Engine
type Engine struct {
	*RouterGroup
	groups []*RouterGroup
	router *router
}

func NewEngine() *Engine {
	e := &Engine{
		groups: make([]*RouterGroup, 0),
		router: newRouter(),
	}
	e.RouterGroup = &RouterGroup{
		prefix: "",
		parent: nil,
		engine: e,
	}
	return e
}

func (e *Engine) Run(addr string) error {
	log.Printf("Listen to %s \n", addr)
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	middlewares := make([]HandlerFunc, 0)
	for _, group := range e.groups {
		if strings.HasPrefix(c.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c.MiddleWares = middlewares
	log.Printf("[%s] %s", c.Method, c.Path)
	e.router.handle(c)
}

func (e *Engine) addGroup(group *RouterGroup) {
	e.groups = append(e.groups, group)
}

// RouterGroup
type RouterGroup struct {
	prefix      string //route prefix
	parent      *RouterGroup
	engine      *Engine
	middlewares []HandlerFunc
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	prefix = group.prefix + prefix
	g := &RouterGroup{
		prefix: prefix,
		parent: group,
		engine: group.engine,
	}
	g.engine.addGroup(g)
	return g
}

func (group *RouterGroup) GET(pattern string, handle HandlerFunc) {
	group.AddRoute(string(GET), pattern, handle)
}

func (group *RouterGroup) POST(pattern string, handle HandlerFunc) {
	group.AddRoute(string(POST), pattern, handle)
}

func (group *RouterGroup) AddRoute(method, pattern string, handle HandlerFunc) {
	method = strings.ToUpper(method)
	group.engine.router.addRouter(MethodType(method), group.prefix+pattern, handle)
}

func (group *RouterGroup) UseMiddleWare(handle HandlerFunc) {
	group.middlewares = append(group.middlewares, handle)
}
