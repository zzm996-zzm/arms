package command

import (
	"fmt"

	"github.com/zzm996-zzm/arms/framework/cobra"
	"github.com/zzm996-zzm/arms/framework/contract"
)

func initEnvCommand() *cobra.Command {
	return envCommand
}

// envCommand 获取当前的 App 环境
var envCommand = &cobra.Command{
	Use:   "env",
	Short: "获取当前的 App 环境",
	Run: func(c *cobra.Command, args []string) {
		// 获取 env 环境
		container := c.GetContainer()
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		// 打印环境
		fmt.Println("environment:", envService.AppEnv())
	},
}
