package demo

import (
	"fmt"

	demoService "github.com/arms/app/provider/demo"
	"github.com/arms/framework/gin"
)

func Register(r *gin.Engine) error {
	// api := NewDemoApi()
	r.Bind(&demoService.DemoServiceProvider{})
	r.GET("/demo", subjuect)
	return nil
}

func subjuect(ctx *gin.Context) {
	fmt.Println("subject")
}
