package http

import (
	"github.com/arms/app/http/module/demo"
	"github.com/arms/framework/gin"
	"github.com/arms/framework/middleware/static"
)

func Routes(r *gin.Engine) {

	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	//TODO:原理
	r.Use(static.Serve("/", static.LocalFile("./vue/dist/", false)))
	demo.Register(r)
}
