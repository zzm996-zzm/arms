package framework

//初始化一个新实例，所以服务器容器的创建服务
type NewInstance func(...interface{}) (interface{}, error)

type ServiceProvider interface {
	//Register 在服务器容器中注册了应该实例化服务的方法是否在注册的时候就实例化这个服务需要参考 IsDefer 接口
	Register(Container) NewInstance

	//Boot 在调用实例化服务的时候会调用，可以把一些准备工作：基础配置，初始化参数的操作放在这里
	//如果Boot返回error，整个服务实例化就会失败，返回错误
	Boot(Container) error

	//IsDefer 决定是否在注册的时候实例化这个服务，如果不是注册的时候实例化，那么就是在第一次 make 的时候进行实例化操作
	//false 表示不需要延迟实例化，在注册的时候就直接实例化
	//true表示延迟实例化
	IsDefer() bool

	//Params param定义传递给NewInstance 的参数，可以自定义多个，建议将container作为第一个参数
	Params(Container) []interface{}

	//Name 代表了这个服务提供者的凭证
	Name() string
}
