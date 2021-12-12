package distributed

import (
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
)

// ArmsAppProvider 提供App的具体实现方法
type LocalDistributedProvider struct {
}

// Register 注册ArmsApp方法
func (app *LocalDistributedProvider) Register(container framework.Container) framework.NewInstance {
	return NewLocalDistributedService
}

// Boot 启动调用
func (app *LocalDistributedProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (app *LocalDistributedProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (app *LocalDistributedProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

// Name 获取字符串凭证
func (app *LocalDistributedProvider) Name() string {
	return contract.DistributedKey
}
