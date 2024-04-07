package geeweb

import (
	"net/http"
	"strings"
)

type HandleFunc func(c *Context)

type MethodType string

const (
	GET  MethodType = "GET"
	POST MethodType = "POST"
)

// Engine
type Engine struct {
	*RouterGroup
	groups []*RouterGroup
}

func NewEngine() *Engine {
	e := &Engine{
		groups: make([]*RouterGroup, 0),
	}
	e.RouterGroup = &RouterGroup{
		prefix: "",
		parent: nil,
		engine: e,
	}
	return e
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (e *Engine) addGroup(group *RouterGroup) {
	e.groups = append(e.groups, group)
}

// RouterGroup
type RouterGroup struct {
	prefix string //route prefix
	parent *RouterGroup
	engine *Engine
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

func (group *RouterGroup) GET(pattern string, handle HandleFunc) {

}

func (group *RouterGroup) POST(pattern string, handle HandleFunc) {

}

func (group *RouterGroup) AddRoute(method, pattern string, handle HandleFunc) {
	method = strings.ToUpper(method)

}
