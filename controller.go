package main

import (
	"context"
	"fmt"
	"log"
	"simple/framework"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)
	fmt.Printf("当前的request 地址是 :%p \n", c.GetRequest())
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(12*time.Second))
	defer cancel()

	// mu := sync.Mutex{}
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// Do real action
		time.Sleep(10 * time.Second)
		c.Json(200, "ok")

		finish <- struct{}{}
	}()
	select {
	//异常退出
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}
	return nil
}

func UserLoginController(c *framework.Context) error {
	// 打印控制器名字
	c.Json(200, "ok, UserLoginController")
	return nil
}

func SubjectGetController(c *framework.Context) error {
	// 打印控制器名字
	c.Json(200, "ok, SubjectGetController")
	return nil
}

func SubjectListController(c *framework.Context) error {
	// 打印控制器名字
	c.Json(200, "ok, SubjectListController")
	return nil
}
