package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arms/framework/gin"
	"github.com/arms/framework/middleware"
	"github.com/arms/framework/provider/app"

	hadeHttp "github.com/arms/app/http"
)

func main() {
	core := gin.New()

	core.Bind(&app.ArmsAppProvider{})
	core.Use(gin.Recovery())
	core.Use(middleware.Cost())

	hadeHttp.Routes(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}

	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
