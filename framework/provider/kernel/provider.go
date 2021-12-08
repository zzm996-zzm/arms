package kernel

import (
	"github.com/arms/framework"
	"github.com/arms/framework/contract"
	"github.com/arms/framework/gin"
)

type ArmsKernelProvider struct {
	HttpEngine *gin.Engine
}

func (provider *ArmsKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewArmsKernelService
}

//Boot 启动的时候判断是否外界注入了Engine 如果注入的话，用注入的。如果没有重新实例化
func (provider *ArmsKernelProvider) Boot(c framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}

	provider.HttpEngine.SetContainer(c)
	return nil
}

func (provider *ArmsKernelProvider) IsDefer() bool {
	return false
}

func (provider *ArmsKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

func (provider *ArmsKernelProvider) Name() string {
	return contract.KernelKey
}
