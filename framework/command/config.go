package command

import (
	"errors"
	"fmt"

	"github.com/kr/pretty"
	"github.com/zzm996-zzm/arms/framework/cobra"
	"github.com/zzm996-zzm/arms/framework/contract"
)

// initConfigCommand 获取配置相关的命令
func initConfigCommand() *cobra.Command {
	configCommand.AddCommand(configGetCommand)
	return configCommand
}

// envCommand 获取当前的App环境
var configCommand = &cobra.Command{
	Use:   "config",
	Short: "获取配置相关信息",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var configGetCommand = &cobra.Command{
	Use:   "get",
	Short: "获取配置项",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		config := container.MustMake(contract.ConfigKey).(contract.Config)
		if len(args) < 1 {
			return errors.New("param cannot be empty")
		}
		configPath := args[0]
		val := config.Get(configPath)
		if val == nil {
			fmt.Println("配置路径 ", configPath, " 不存在")
			return nil
		}

		fmt.Printf("%# v\n", pretty.Formatter(val))
		return nil

	},
}
