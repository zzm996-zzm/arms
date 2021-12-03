package main

import (
	"fmt"
	"simple/framework"
	"time"
)

func UserLoginController(c *framework.Context) error {
	foo, _ := c.QueryString("foo", "def")
	// 等待10s才结束执行
	time.Sleep(10 * time.Second)
	// 输出结果
	c.SetOkStatus().Json("ok, UserLoginController: " + foo)
	return nil
}

func UserLoginController2(c *framework.Context) error {
	// 打印控制器名字
	c.Json("ok, UserLoginController222222222")
	return nil
}

func SubjectGetController(c *framework.Context) error {
	// 打印控制器名字
	fmt.Println(c.QueryAll())
	c.Json("ok, SubjectGetController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	// 打印控制器名字
	c.Json("ok, SubjectListController")
	return nil
}
