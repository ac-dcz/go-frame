package geeweb

import "sync"

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
	if height == len(parts) || root.part[0] == '*' {
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
	roots       map[string]*tire
	handleFuncs map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		mu:          sync.Mutex{},
		roots:       make(map[string]*tire),
		handleFuncs: make(map[string]HandleFunc),
	}
}

func (r *router) addRouter(method MethodType, pattern string, handle HandleFunc) {
	// key := string(method) + "-" + pattern

}
