package app

import (
	"github.com/arms/framework"
	"github.com/arms/framework/contract"
)

// ArmsAppProvider 提供App的具体实现方法
type ArmsAppProvider struct {
	BaseFolder string
}

// Register 注册ArmsApp方法
func (app *ArmsAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewArmsApp
}

// Boot 启动调用
func (app *ArmsAppProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (app *ArmsAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (app *ArmsAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, app.BaseFolder}
}

// Name 获取字符串凭证
func (app *ArmsAppProvider) Name() string {
	return contract.AppKey
}
