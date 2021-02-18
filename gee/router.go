package gee

import "net/http"

// Gee路由管理

type router struct {
	handlers map[string]HandlerFunc // 静态路由存储
}

// 创建Gee router对象
func newRouter() *router {
	return &router{make(map[string]HandlerFunc)}
}

// 为请求存储相应的处理函数
func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	router.handlers[key] = handler
}

// 根据请求查询相应的处理函数并执行
func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := router.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
