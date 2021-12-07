package main

import (
	"time"

	"github.com/arms/framework/gin"
)

func UserLoginController(c *gin.Context) {

	// foo, _ := c.DefaultQueryString("foo", "def")
	// c.ISetOkStatus().IJson("ok, UserLoginController: " + foo)

	// durationInMillis := durationInSeconds * 1000
	endTimeInMillis := time.Now().UnixNano() / 1e6
	c.ISetOkStatus().IJson(endTimeInMillis)
	// startTimeInMillis := endTimeInMillis - durationInMillis
}

func UserRegisterController(c *gin.Context) {
	// start := float64(time.Now().Unix())
	// metrics.recordTimestamp("register", start)
	foo, _ := c.DefaultQueryString("foo", "def")
	// 输出结果
	c.ISetOkStatus().IJson("ok, UserRegisterController: " + foo)
	// metrics.recordResponseTime("register", float64(time.Now().Unix())-start)
}
