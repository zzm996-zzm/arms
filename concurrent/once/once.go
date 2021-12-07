package main

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	sync.Mutex
}

func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 0 {
		return o.doSlow(f)
	}

	return nil
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.done) == 1
}

func (o *Once) doSlow(f func() error) error {
	o.Lock()
	defer o.Unlock()
	var err error
	//双重检查
	if o.done == 0 {
		err = f()
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}

	return err
}
