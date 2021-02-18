package gee

import "net/http"

// Gee路由管理

type HandlerFunc func(http.ResponseWriter, *http.Request)

type router struct {
	handlers map[string]HandlerFunc
}

// New is the constructor of gee.Engine
func newRouter() *router {
	return &router{make(map[string]HandlerFunc)}
}

func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	router.handlers[key] = handler
}

func (router *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := router.handlers[key]; ok {
		handler(c.Writer, c.Req)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
