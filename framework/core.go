package framework

import (
	"log"
	"net/http"
	"strings"
)

// 框架核心结构

type Core struct {
	router      map[string]*Tree // all routers
	middlewares []ControllerHandler
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

func (c *Core) addRouter(method, url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router[method].AddRouter(url, allHandlers...); err != nil {
		log.Fatal("add router error: ", err)
	}
}
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	c.addRouter("GET", url, handlers...)
}
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	c.addRouter("POST", url, handlers...)
}
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	c.addRouter("PUT", url, handlers...)
}
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	c.addRouter("DELETE", url, handlers...)
}

// 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) *node {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

func (c *Core) Use(middleware ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middleware...)
}

// 框架核心结构实现Handler接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// 封装自定义context
	ctx := NewContext(request, response)

	// 寻找路由
	node := c.FindRouteByRequest(request)
	if node == nil {
		// 如果没有找到，这里打印日志
		ctx.Json("not found").SetStatus(404)
		return
	}

	// 设置路由参数
	params := node.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	ctx.SetHandlers(node.handlers)

	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	if err := ctx.Next(); err != nil {
		ctx.Json("inner error").SetStatus(404)
		return
	}
}
