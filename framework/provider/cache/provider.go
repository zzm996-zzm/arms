package cache

import (
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/provider/cache/services"
	"strings"
)


// CacheProvider 提供App的具体实现方法
type CacheProvider struct {
	framework.ServiceProvider
}

// Register 注册方法
func (h *CacheProvider) Register(container framework.Container) framework.NewInstance {
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	driver := strings.ToLower(configService.GetString("cache.driver"))

	// 根据driver的配置项确定
	switch driver {
	case "redis":
		return services.NewRedisCache
	case "memory":
		return nil
	default:
		return nil
	}

}

// Boot 启动调用
func (h *CacheProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (h *CacheProvider) IsDefer() bool {
	return true
}

// Params 获取初始化参数
func (h *CacheProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

// Name 获取字符串凭证
func (h *CacheProvider) Name() string {
	return contract.CacheKey
}
