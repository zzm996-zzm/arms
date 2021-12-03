package main

import (
	"simple/framework"
)

// 注册路由规则
func registerRouter(core *framework.Core) {
	// 在core中使用middleware.Test3() 为单个路由增加中间件
	core.Get("/:id/cc/aa", SubjectGetController) // c.middlewares len 3 cap 4
	core.Get("/:id/aa/bb", UserLoginController)  //

	// 批量通用前缀
	// subjectApi := core.Group("/subject")
	// {
	// 	// 在group中使用middleware.Test3() 为单个路由增加中间件
	// 	// subjectApi.Get("/:id", middleware.Test3(), SubjectGetController)
	// }

}
