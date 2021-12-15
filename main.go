package main

import (
	"github.com/zzm996-zzm/arms/app/console"
	armsHttp "github.com/zzm996-zzm/arms/app/http"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/provider/app"
	"github.com/zzm996-zzm/arms/framework/provider/cache"
	"github.com/zzm996-zzm/arms/framework/provider/config"
	"github.com/zzm996-zzm/arms/framework/provider/distributed"
	"github.com/zzm996-zzm/arms/framework/provider/env"
	"github.com/zzm996-zzm/arms/framework/provider/kernel"
	"github.com/zzm996-zzm/arms/framework/provider/log"
	"github.com/zzm996-zzm/arms/framework/provider/orm"
)

func main() {
	// 初始化服务容器
	container := framework.NewAppContainer()
	// 绑定App服务提供者 有顺序要求 前3个强依赖关系
	container.Bind(&app.ArmsAppProvider{})
	container.Bind(&env.EnvProvider{})
	container.Bind(&config.ConfigProvider{})
	container.Bind(&log.LogProvider{})
	// 后续初始化需要绑定的服务提供者...
	container.Bind(&distributed.LocalDistributedProvider{})
	container.Bind(&orm.GormProvider{})
	container.Bind(&cache.CacheProvider{})

	// 将HTTP引擎初始化,并且作为服务提供者绑定到服务容器中
	if engine, err := armsHttp.NewHttpEngine(); err == nil {
		err := container.Bind(&kernel.ArmsKernelProvider{HttpEngine: engine})
		if err != nil {
			return
		}
	}

	// 运行root命令
	err := console.RunCommand(container)
	if err != nil {
		return 
	}
}
