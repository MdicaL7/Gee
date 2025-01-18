package gee

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots 存储不同请求方式的 trie 树根节点，handlers 存储不同请求方式的 HandlerFunc。
// roots 类型是 map[string]*node，key 是请求方式，value 是 *node，即路由树的根节点。
// handlers 类型是 map[string]HandlerFunc，key 是请求方式，value 是 HandlerFunc，即路由处理函数。

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 函数用来解析路由规则，将路由规则转为 []string。例如 /p/:lang 转换为 [p, :lang]，/p/*filepath 转换为 [p, *filepath]。
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { // 如果是通配符，则直接break
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 检查路由冲突,如果已经存在该路由，则直接panic
	// 例如，如果已经存在 GET - /hello/:name 这个路由规则，再添加 GET - /hello/geektutu 这个路由规则，就会触发 panic。
	if _, existing := r.getRoute(method, pattern); existing != nil {
		panic(fmt.Sprintf("route conflict: %s %s", method, pattern))
	}
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	// 如果不存在该请求方式的路由树，则新建一个
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 插入路由规则
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts { //遍历路由的每一部分
			if part[0] == ':' { //如果是:开头的，则是参数
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// handle 函数用来处理请求，从请求方法中找到路由规则，解析路由规则，找到对应的 HandlerFunc，执行 HandlerFunc。
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
		c.String(404, "404 NOT FOUND: %s\n", c.Path)
	}
}
