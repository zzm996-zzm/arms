package console

import (
	"github.com/arms/framework"
	"github.com/arms/framework/cobra"
	"github.com/arms/framework/command"
)

func RunCommand(container framework.Container) error {
	//根Command
	var rootCmd = &cobra.Command{
		Use:   "arms",
		Short: "arms 命令",
		Long:  "arms 框架提供的命令行工具，使用这个命令行工具可以方便执行框架自带命令,也能很方便的编写业务命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	rootCmd.SetContainer(container)
	command.AddKernelCommand(rootCmd)
	AddAppCommand(rootCmd)
	return rootCmd.Execute()
}

// 绑定业务的命令
func AddAppCommand(rootCmd *cobra.Command) {
	//  demo 例子
	// rootCmd.AddCommand(demo.InitFoo())
	// rootCmd.AddCronCommand("* * * * * *", demo.Foo1Command)
	// rootCmd.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}
