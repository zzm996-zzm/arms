package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"strconv"
	"time"
)

func GetBaseConfig(c framework.Container) *contract.RedisConfig{
	logger := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.RedisConfig{Options:&redis.Options{}}
	opt := WithConfigPath("redis")
	err := opt(c,config)
	if err != nil {
		logger.Error(context.Background(),"parse cache config error", nil)
		return nil
	}
	return config
}

func WithConfigPath(path string) contract.RedisOption {
	return func(container framework.Container,config *contract.RedisConfig) error{
		configService :=container.MustMake(contract.ConfigKey).(contract.Config)
		conf := configService.GetStringMapString(path)

		if host,ok := conf["host"];ok{
			if port,ok1 := conf["port"];ok1{
				config.Addr = host + ":" + port
			}
		}

		if db,ok := conf["db"];ok{
			t,err := strconv.Atoi(db)
			if err != nil {
				return err
			}
			config.DB = t
		}

		if username, ok := conf["username"]; ok {
			config.Username = username
		}

		if password, ok := conf["password"]; ok {
			config.Password = password
		}

		if timeout, ok := conf["timeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.DialTimeout = t
		}


		if timeout, ok := conf["readTimeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.ReadTimeout = t
		}

		if timeout, ok := conf["writeTimeout"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.WriteTimeout = t
		}

		if cnt, ok := conf["connMinIdle"]; ok {
			t, err := strconv.Atoi(cnt)
			if err != nil {
				return err
			}
			config.MinIdleConns = t
		}

		if max, ok := conf["connMaxOpen"]; ok {
			t, err := strconv.Atoi(max)
			if err != nil {
				return err
			}
			config.PoolSize = t
		}

		if timeout, ok := conf["connMaxLifetime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.MaxConnAge = t
		}

		if timeout, ok := conf["connMaxIdleTime"]; ok {
			t, err := time.ParseDuration(timeout)
			if err != nil {
				return err
			}
			config.IdleTimeout = t
		}

		return nil
 	}
}

// WithRedisConfig 表示自行配置redis的配置信息
func WithRedisConfig(f func(options *contract.RedisConfig)) contract.RedisOption {
	return func(container framework.Container, config *contract.RedisConfig) error {
		f(config)
		return nil
	}
}