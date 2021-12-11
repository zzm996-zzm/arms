package command

import "github.com/arms/framework/cobra"

func AddKernelCommand(root *cobra.Command) {
	root.AddCommand(initAppCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(initEnvCommand())
	root.AddCommand(initConfigCommand())
	root.AddCommand(initBuildCommand())
	root.AddCommand(initDevCommand())
	root.AddCommand(initProviderCommand())
}
