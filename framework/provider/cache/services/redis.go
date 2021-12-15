package services

import (
	"context"
	"errors"
	redisv8 "github.com/go-redis/redis/v8"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/provider/redis"
	"sync"
	"time"
)


type RedisCache struct{
	container framework.Container
	client    *redisv8.Client
	lock      *sync.RWMutex
}



var _ contract.Cache = new(RedisCache)

func NewRedisCache(params ...interface{}) (interface{},error){
	container := params[0].(framework.Container)
	if !container.IsBind(contract.RedisKey){
		err := container.Bind(&redis.RedisProvider{})
		if err !=nil {
			return nil,err
		}
	}

	redisService := container.MustMake(contract.RedisKey).(contract.Redis)
	client ,err := redisService.GetClient(redis.WithConfigPath("cache"))
	if err != nil {
		return nil,err
	}

	//返回RedisCache实例

	redisCache:= &RedisCache{
		container: container,
		client: client,
		lock: &sync.RWMutex{},
	}

	return redisCache,nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val,err := r.client.Get(ctx,key).Result()
	if errors.Is(err,redisv8.Nil){
		return val,ErrKeyNotFound
	}
	return val,err

}

func (r *RedisCache) GetObj(ctx context.Context, key string,model interface{})error{
	cmd := r.client.Get(ctx,key)
	if errors.Is(cmd.Err(),redisv8.Nil){
		return ErrKeyNotFound
	}
	err := cmd.Scan(model)
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := r.client.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redisv8.StringCmd,0,len(keys))

	for _,key := range keys{
		cmds = append(cmds,pipeline.Get(ctx,key))
	}
	_,err := pipeline.Exec(ctx)
	if err != nil {
		return nil,err
	}
	errs := make([]string,0,len(keys))
	for _,cmd := range cmds{
		val,err := cmd.Result()
		if err != nil {
			errs = append(errs,err.Error())
			continue
		}
		key := cmd.Args()[1].(string)
		vals[key] = val
	}
	return vals,nil
}

func (r *RedisCache) Set(ctx context.Context,key,val string,timeout time.Duration) error{
	return r.client.Set(ctx,key,val,timeout).Err()
}

func (r *RedisCache) SetForever(ctx context.Context, key string, val string) error {
	return r.Set(ctx,key,val,NoneDuration)
}

// SetForeverObj 设置某个key和对象到缓存，不带超时时间，对象必须实现 https://pkg.go.dev/encoding#BinaryMarshaler
func (r *RedisCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return r.SetObj(ctx, key, val, NoneDuration)
}

func (r *RedisCache) SetObj(ctx context.Context, key string,val interface{},timeout time.Duration) error{
	return r.client.Set(ctx,key,val,timeout).Err()
}

func (r *RedisCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipline := r.client.Pipeline()
	cmds := make([]*redisv8.StatusCmd,0,len(data))
	for k,v:= range data{
		cmds = append(cmds,pipline.Set(ctx,k,v,timeout))
	}
	_,err := pipline.Exec(ctx)
	return err
}

// SetTTL 设置某个key的超时时间
func (r *RedisCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return r.client.Expire(ctx, key, timeout).Err()
}

// GetTTL 获取某个key的超时时间
func (r *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()


}

func (r *RedisCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return r.client.IncrBy(ctx, key, step).Result()
}

func (r *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, 1).Result()
}

func (r *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return r.client.IncrBy(ctx, key, -1).Result()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) DelMany(ctx context.Context, keys []string) error {
	pipline := r.client.Pipeline()
	cmds := make([]*redisv8.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipline.Del(ctx, key))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (r *RedisCache) Remember(ctx context.Context,key string,timeout time.Duration,rememberFunc contract.RememberFunc, obj interface{}) error{
	err := r.GetObj(ctx,key,obj)
	//fast path
	if err == nil {
		return nil
	}
	if !errors.Is(err,ErrKeyNotFound){
		return err
	}

	objNew,err := rememberFunc(ctx,r.container)
	if err != nil {
		return err
	}

	if err := r.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}
	if err := r.GetObj(ctx, key, obj); err != nil {
		return err
	}
	return nil
}