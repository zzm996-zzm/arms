package framework

import (
	"context"
	"log"
	"time"
)

func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	//使用函数回调
	return func(c *Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		//执行业务逻辑之前的预操作： 初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.request.WithContext(durationCtx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			fun(c)

			finish <- struct{}{}
		}()

		select {
		//异常退出
		case p := <-panicChan:
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			log.Println(p)
			c.Json("panic").SetStatus(500)
		case <-finish:
		case <-durationCtx.Done():
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			c.Json("time out").SetStatus(500)
			c.SetHasTimeout()
		}
		return nil
	}
}
