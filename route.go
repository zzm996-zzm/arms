package main

import (
	"github.com/arms/framework/gin"
)

// 注册路由规则
func registerRouter(core *gin.Engine) {
	// 静态路由+HTTP方法匹配
	core.GET("/user/login", UserLoginController)
	core.GET("/user/register", UserRegisterController)
	core.GET("/subject", SubjectListController)
}
