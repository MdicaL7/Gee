package gee

import (
	"net/http"
)

// HandlerFunc 定义 gee 使用的请求处理程序
type HandlerFunc func(c *Context)

// Engine 实现handler接口方法ServerHTTP
type Engine struct {
	router *router //路由映射表，请求方法和路由地址不同映射至不同的处理方法
}

// New gee.Engine的构造函数
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET 添加GET请求函数
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST 添加POST请求函数
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run 启动http服务的方法
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 实现的handle接口中ServeHTTP方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
