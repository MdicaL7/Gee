package gee

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle 方法的实现很简单，从请求中取到 Method 和 Path，然后拼接成 key，从 handlers 中取出对应的 HandlerFunc，如果取不到说明路由未注册，返回 404 NOT FOUND。
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
		c.String(404, "404 NOT FOUND: %s\n", c.Path)
	}
}
