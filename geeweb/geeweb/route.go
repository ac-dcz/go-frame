package geeweb

import (
	"net/http"
	"strings"
	"sync"
)

/*
 *path:
	1. /a/b/c
	2. /a/:b/c
	3. /a/*b
*/

// tireNode
type tireNode struct {
	pattern  string // 初试 匹配模式
	part     string
	isWild   bool // 是否支持泛型匹配
	children []*tireNode
}

func (node *tireNode) findOneChild(part string) *tireNode {
	for _, child := range node.children {
		if !child.isWild && child.part == part {
			return child
		}
	}
	return nil
}

func (node *tireNode) findAllChild(part string) []*tireNode {
	var res []*tireNode
	for _, child := range node.children {
		if child.isWild || child.part == part {
			res = append(res, child)
		}
	}
	return res
}

type tire struct {
	root *tireNode
}

func newTire() *tire {
	t := &tire{}
	t.root = &tireNode{
		pattern:  "",
		part:     "",
		isWild:   false,
		children: make([]*tireNode, 0),
	}
	return t
}

func (tr *tire) insertRoute(pattern string, parts []string) *tireNode {
	return _insert(tr.root, pattern, parts, 0)
}

func _insert(root *tireNode, pattern string, parts []string, height int) *tireNode {
	if height == len(parts) {
		root.pattern = pattern
		return root
	}
	child := root.findOneChild(parts[height])
	if child == nil {
		child = &tireNode{
			pattern:  "",
			part:     parts[height],
			isWild:   parts[height][0] == ':' || parts[height][0] == '*',
			children: make([]*tireNode, 0),
		}
		root.children = append(root.children, child)
	}
	return _insert(child, pattern, parts, height+1)
}

func (tr *tire) searchRoute(parts []string) *tireNode {
	return _search(tr.root, parts, 0)
}

func _search(root *tireNode, parts []string, height int) *tireNode {
	if height == len(parts) || (len(root.part) > 0 && root.part[0] == '*') {
		return root
	}
	childs := root.findAllChild(parts[height])
	for _, child := range childs {
		if node := _search(child, parts, height+1); node != nil {
			return node
		}
	}
	return nil
}

type router struct {
	mu          sync.Mutex
	roots       map[MethodType]*tire
	handleFuncs map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		mu:          sync.Mutex{},
		roots:       make(map[MethodType]*tire),
		handleFuncs: make(map[string]HandleFunc),
	}
}

func (r *router) addRouter(method MethodType, pattern string, handle HandleFunc) {
	key := string(method) + "-" + pattern
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handleFuncs[key] = handle
	root := r.roots[method]
	if root == nil {
		root = newTire()
		r.roots[method] = root
	}
	root.insertRoute(pattern, parsePattern(pattern))
}

func parsePattern(pattern string) []string {
	items := strings.Split(pattern, "/")
	var parts []string
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) getRouter(method MethodType, pattern string) *tireNode {
	items := parsePattern(pattern)
	r.mu.Lock()
	defer r.mu.Unlock()
	root := r.roots[method]
	if root == nil {
		return nil
	}
	return root.searchRoute(items)
}

func parseURL(rawPattern, pattern string) map[string]string {
	item1, item2 := parsePattern(rawPattern), parsePattern(pattern)
	params := make(map[string]string)
	for i, item := range item1 {
		if item[0] == ':' {
			params[item[1:]] = item2[i]
		} else if item[0] == '*' {
			params[item[1:]] = strings.Join(item2[i:], "/")
			break
		}
	}
	return params
}

func (r *router) handle(c *Context) {
	if node := r.getRouter(c.Method, c.Path); node != nil {
		c.Params = parseURL(node.pattern, c.Path)
		key := string(c.Method) + "-" + node.pattern
		r.handleFuncs[key](c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found: %s\n", c.Path)
	}
}
