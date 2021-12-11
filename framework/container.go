package framework

import (
	"errors"
	"sync"
)

type Container interface {
	//Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回error
	Bind(provider ServiceProvider) error

	//判断凭证是否已经被绑定
	IsBind(key string) bool

	//Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)

	//MustMake 根据关键字获取一个服务，如果这个关键字凭证为绑定服务提供者，那么会panic
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}

	//MakeNew根据关键字凭证获取一个服务，只是这个服务并不是单例的
	//是根据服务提供者注册的启动函数和传递的param实例化出来的
	//这个函数在需要为不同参数启动不同实例的时候非常有用

	MakeNew(key string, params []interface{}) (interface{}, error)
}

type AppContainer struct {
	Container //强制要求实现Container接口

	//存储注册的服务提供者，key为字符串凭证
	providers map[string]ServiceProvider

	//instance 存储具体的实例，key为字符串凭证
	instances map[string]interface{}

	lock sync.RWMutex
}

func NewAppContainer() *AppContainer {
	return &AppContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

func (app *AppContainer) NameList() []string {
	ret := []string{}
	for _, provider := range app.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}

func (app *AppContainer) Bind(provider ServiceProvider) error {
	//写锁
	app.lock.Lock()
	key := provider.Name()
	app.providers[key] = provider
	app.lock.Unlock()

	if !provider.IsDefer() {
		return app.instanceBind(provider)
	}
	return nil
}

func (app *AppContainer) instanceBind(provider ServiceProvider) error {
	instance, err := app.newInstance(provider, nil)
	if err != nil {
		return errors.New(err.Error())
	}
	app.lock.Lock()
	defer app.lock.Unlock()
	app.instances[provider.Name()] = instance

	return nil
}

func (app *AppContainer) IsBind(key string) bool {
	return app.findServiceProvider(key) != nil
}

func (app *AppContainer) findServiceProvider(key string) ServiceProvider {
	app.lock.RLock()
	defer app.lock.RUnlock()
	if sp, ok := app.providers[key]; ok {
		return sp
	}
	return nil

}

func (app *AppContainer) findInstance(key string) interface{} {
	app.lock.RLock()
	defer app.lock.RUnlock()
	if sp, ok := app.instances[key]; ok {
		return sp
	}
	return nil
}

// Make 方式调用内部的 make 实现
func (app *AppContainer) Make(key string) (interface{}, error) {
	return app.make(key, nil, false)
}

func (app *AppContainer) MustMake(key string) interface{} {
	serv, err := app.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

// MakeNew 方式使用内部的 make 初始化
func (app *AppContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return app.make(key, params, true)
}

func (app *AppContainer) newInstance(provider ServiceProvider, params []interface{}) (interface{}, error) {
	if err := provider.Boot(app); err != nil {
		return nil, err
	}

	if params == nil {
		params = provider.Params(app)
	}

	method := provider.Register(app)
	instance, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return instance, err
}

// 真正的实例化一个服务
//TODO:拆分成更单一的函数
func (app *AppContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	// 查询是否已经注册了这个服务提供者，如果没有注册，则返回错误
	sp := app.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return app.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	ins := app.findInstance(key)
	if ins != nil {
		return ins, nil
	}

	//绑定进instances数组
	err := app.instanceBind(sp)
	if err != nil {
		return nil, errors.New("new " + key + "instance error")
	}

	ins = app.findInstance(key)
	//再次判断确保稳定性
	if ins == nil {
		return nil, errors.New("the instance is not initialized")
	}

	return ins, nil

}
