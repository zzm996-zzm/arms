package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simple/framework"
	"syscall"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)
	server := &http.Server{
		// 自定义的请求核心处理函数
		Handler: core,
		// 请求监听地址
		Addr: ":8080",
	}

	go func() {
		server.ListenAndServe()
	}()

	//当前goroutine 等待的信号量
	quit := make(chan os.Signal)

	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit

	//触发这些信号则关闭
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
