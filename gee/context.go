package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// Context 存储一次请求的上下文信息
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

/**
构建Context对象
*/

func newContext(writer http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: writer,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

/**
获取POST请求中某一字段的值
*/
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

/**
获取GET请求中某一字段的值
*/
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

/**
为响应设置状态码信息
*/
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

/**
为响应设置头信息
*/
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

/**
快速构造String响应
*/
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

/**
快速构造JSON响应
*/
func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

/**
快速构造Data响应
*/
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

/**
快速构造HTML响应
*/
func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}
