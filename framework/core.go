package framework

import (
	"log"
	"net/http"
	"strings"
)

// 框架核心结构

type Core struct {
	router map[string]*Tree // all routers
}

// 初始化Core结构
func NewCore() *Core {
	// 初始化路由
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

// ==== http method wrap end

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

func (c *Core) addRouter(method, url string, handler ControllerHandler) {
	if err := c.router[method].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}
func (c *Core) Get(url string, handler ControllerHandler) {
	c.addRouter("GET", url, handler)
}
func (c *Core) Post(url string, handler ControllerHandler) {
	c.addRouter("POST", url, handler)
}
func (c *Core) Put(url string, handler ControllerHandler) {
	c.addRouter("PUT", url, handler)
}
func (c *Core) Delete(url string, handler ControllerHandler) {
	c.addRouter("DELETE", url, handler)
}

// 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) ControllerHandler {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

// 框架核心结构实现Handler接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// 封装自定义context
	ctx := NewContext(request, response)

	// 寻找路由
	router := c.FindRouteByRequest(request)
	if router == nil {
		// 如果没有找到，这里打印日志
		ctx.Json(404, "not found")
		return
	}

	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	if err := router(ctx); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
