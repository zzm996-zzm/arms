package demoService

import (
	"github.com/zzm996-zzm/arms/framework"
)

// 具体的接口实例
type DemoService struct {
	// 实现接口
	Service

	// 参数
	c framework.Container
}

// 初始化实例的方法
func NewDemoService(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)

	// 返回实例
	return &DemoService{c: c}, nil
}

// 实现接口
func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}
