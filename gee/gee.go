package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc 定义 gee 使用的请求处理程序
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// Engine 实现handler接口方法ServerHTTP
type Engine struct {
	router map[string]HandlerFunc //路由映射表，请求方法和路由地址不同映射至不同的处理方法
}

// New gee.Engine的构造函数
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern //key 由请求方法和静态路由地址构成,如GET-/、GET-/hello、POST-/hello
	engine.router[key] = handler
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
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 page not found:%s\n", req.URL)
	}
}
