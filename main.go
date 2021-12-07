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
	"github.com/arms/metrics"
	"github.com/arms/provider/demo"
)

func main() {
	core := gin.New()

	core.Bind(&demo.DemoServiceProvider{})
	core.Use(gin.Recovery())
	core.Use(middleware.Cost())

	registerRouter(core)
	server := &http.Server{
		Handler: core,
		Addr:    ":8080",
	}

	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	console := metrics.NewConsoleReporter(metrics.NewRedisMetricsStorage())
	console.StartRepeatedReport(60)

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
