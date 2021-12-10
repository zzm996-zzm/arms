package config

import (
	"path/filepath"

	"github.com/arms/framework"
	"github.com/arms/framework/contract"
)

type ConfigProvider struct {
}

// Register registe a new function for make a service instance
func (provider *ConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewConfigService
}

// Boot will called when the service instantiate
func (provider *ConfigProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *ConfigProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *ConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.ArmsApp)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

// Name / Name define the name for this service
func (provider *ConfigProvider) Name() string {
	return contract.ConfigKey
}
