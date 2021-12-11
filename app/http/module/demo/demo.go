package demo

import (
	demoService "github.com/arms/app/provider/demo"
	"github.com/arms/framework/gin"
)

func Register(r *gin.Engine) error {
	// api := NewDemoApi()
	r.Bind(&demoService.DemoServiceProvider{})
	r.GET("/demo", subjuect)
	return nil
}

func subjuect(c *gin.Context) {
	// 获取password
	// configService := c.MustMake(contract.ConfigKey).(contract.Config)
	// password := configService.GetString("database.mysql.password")

	// 打印出来
	c.JSON(200, "test success")
}
