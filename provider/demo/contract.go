package demo

const Key = "arms:demo"

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}
