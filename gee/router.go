package gee

import (
	"net/http"
	"strings"
)

// Gee路由管理
// 引入动态路由

type router struct {
	// handlers map[string]HandlerFunc // 静态路由存储
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// 创建Gee router对象
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析请求，生成parts
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	res := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			res = append(res, item)
			// 路由中出现一个*即可
			if item[0] == '*' {
				break
			}
		}
	}
	return res
}

// 为请求存储相应的处理函数
func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := router.roots[method]; !ok {
		router.roots[method] = &node{} // 前缀树的根节点是空节点
	}
	router.roots[method].insert(pattern, parts, 0)
	router.handlers[key] = handler
}

// 路由匹配(前缀树匹配代替哈希匹配)
func (router *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string) // 动态路由匹配的参数

	root, ok := router.roots[method]
	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)

	if node != nil {
		parts := parsePattern(node.pattern)
		// 匹配参数
		for index, item := range parts {
			if item[0] == ':' {
				params[item[1:]] = searchParts[index]
			}
			if item[0] == '*' && len(item) > 1 {
				params[item[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

// 根据请求查询相应的处理函数并执行
func (router *router) handle(c *Context) {
	n, params := router.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + c.Path
		// 执行路由处理函数
		router.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
