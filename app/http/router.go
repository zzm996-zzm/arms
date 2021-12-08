package http

import (
	"github.com/arms/app/http/module/demo"
	"github.com/arms/framework/gin"
)

func Routes(r *gin.Engine) {

	r.Static("/dist/", "./dist/")
	demo.Register(r)
}
