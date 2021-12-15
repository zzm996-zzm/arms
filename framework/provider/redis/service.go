package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"sync"
)

type Redis struct{
	container framework.Container
	clients  map[string]*redis.Client
	lock  *sync.RWMutex
}

func NewRedis(params ...interface{})(interface{},error){
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}
	return &Redis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

func (r *Redis)GetClient(options ...contract.RedisOption)  (*redis.Client, error) {
	container := r.container
	//读取默认配置
	config := GetBaseConfig(container)

	//option修改
	for _, opt := range options {
		if err := opt(container, config); err != nil {
			return nil, err
		}
	}

	//如果最终的config没有设置key,就生成key
	key := config.UniqKey()
	r.lock.RLock()
	if client,ok:=r.clients[key];ok{
		r.lock.RUnlock()
		return client,nil
	}
	r.lock.RUnlock()
	//实例化
	r.lock.Lock()
	defer r.lock.Unlock()
	client := redis.NewClient(config.Options)
	//挂在到map中
	r.clients[key] = client


	return client,nil
}