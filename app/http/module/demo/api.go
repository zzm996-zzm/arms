package demo

import (
	demoService "github.com/zzm996-zzm/arms/app/provider/demo"
	"github.com/zzm996-zzm/arms/framework/gin"
)

type DemoApi struct {
}

func NewDemoApi() *DemoApi {
	return &DemoApi{}
}
func Register(r *gin.Engine) error {
	api := NewDemoApi()
	r.Bind(&demoService.DemoServiceProvider{})
	r.GET("/demo", api.DemoOrm)
	return nil
}
