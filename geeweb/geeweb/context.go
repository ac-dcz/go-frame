package geeweb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type H map[string]string

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Method MethodType
	Path   string
	Params map[string]string //URL中的参数
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:      w,
		R:      r,
		Method: MethodType(r.Method),
		Path:   r.URL.Path,
	}
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) PostForm(key string) string {
	return c.R.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) AddHeader(key, val string) {
	c.W.Header().Add(key, val)
}

func (c *Context) SetStatusCode(code int) {
	c.W.WriteHeader(code)
}

func (c *Context) JSON(code int, obj any) {
	c.SetStatusCode(code)
	c.AddHeader("Content-Type", "application/json")
	if err := json.NewEncoder(c.W).Encode(obj); err != nil {
		log.Printf("Context: Write JSON DATA error %v\n", err) //本地日志打印
		http.Error(c.W, err.Error(), 500)                      //服务器内部错误
	}
}

func (c *Context) String(code int, format string, v ...interface{}) {
	c.SetStatusCode(code)
	c.AddHeader("Content-Type", "text/plain")
	fmt.Fprintf(c.W, format, v...)
}

func (c *Context) Data(code int, data []byte) {
	c.SetStatusCode(code)
	c.W.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetStatusCode(code)
	c.AddHeader("Content-Type", "text/html")
	c.W.Write([]byte(html))
}
