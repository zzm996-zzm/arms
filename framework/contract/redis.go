package contract

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zzm996-zzm/arms/framework"
)

const RedisKey = "arms:redis"

// RedisConfig 为hade定义的Redis配置结构
type RedisConfig struct {
	*redis.Options
}

// UniqKey 用来唯一标识一个RedisConfig配置
func (config *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}

type RedisOption func(container framework.Container,config *RedisConfig) error

type Redis interface{
	GetClient(option ...RedisOption) (*redis.Client,error)
}