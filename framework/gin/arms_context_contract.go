package gin

import "github.com/zzm996-zzm/arms/framework/contract"

// MustMakeApp 从容器中获取App服务
func (c *Context) MustMakeApp() contract.ArmsApp {
	return c.MustMake(contract.AppKey).(contract.ArmsApp)
}

// MustMakeKernel 从容器中获取Kernel服务
func (c *Context) MustMakeKernel() contract.Kernel {
	return c.MustMake(contract.KernelKey).(contract.Kernel)
}

// MustMakeConfig 从容器中获取配置服务
func (c *Context) MustMakeConfig() contract.Config {
	return c.MustMake(contract.ConfigKey).(contract.Config)
}

// MustMakeLog 从容器中获取日志服务
func (c *Context) MustMakeLog() contract.Log {
	return c.MustMake(contract.LogKey).(contract.Log)
}

// MustMakeApp 从容器中获取gorm服务
func (c *Context) MustMakeOrm() contract.ORM {
	return c.MustMake(contract.ORMKey).(contract.ORM)
}
