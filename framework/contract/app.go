package contract

const AppKey = "arms:app"

type ArmsApp interface {
	//获取当前服务的唯一AppID 可用于分布式锁
	AppID() string
	//Version 定义当前版本
	Version() string
	//BaseFolder 定义项目基础地址
	BaseFolder() string
	//ConfigFolder 定义了配置文件路径
	ConfigFolder() string
	//LogFolder 定义了日志所在路径
	LogFolder() string
	// ProviderFolder 定义业务自己的服务提供者地址
	ProviderFolder() string
	// MiddlewareFolder 定义业务自己定义的中间件
	MiddlewareFolder() string
	// CommandFolder 定义业务定义的命令
	CommandFolder() string
	// RuntimeFolder 定义业务的运行中间态信息
	RuntimeFolder() string
	// TestFolder 存放测试所需要的信息
	TestFolder() string

	// AppFolder 定义业务代码所在的目录，用于监控文件变更使用
	AppFolder() string

	HttpFolder() string

	LoadAppConfig(v map[string]string)
}
