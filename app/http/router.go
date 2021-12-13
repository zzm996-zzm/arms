package http

import (
	"github.com/zzm996-zzm/arms/app/http/module/demo"
	"github.com/zzm996-zzm/arms/framework/gin"
	ginSwagger "github.com/zzm996-zzm/arms/framework/middleware/gin-swagger"
	"github.com/zzm996-zzm/arms/framework/middleware/gin-swagger/swaggerFiles"
	"github.com/zzm996-zzm/arms/framework/middleware/static"
)

func Routes(r *gin.Engine) {

	// /路径先去./dist目录下查找文件是否存在，找到使用文件服务提供服务
	//TODO:原理
	r.Use(static.Serve("/", static.LocalFile("./vue/dist/", false)))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	demo.Register(r)
}
