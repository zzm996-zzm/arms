package command

import "github.com/zzm996-zzm/arms/framework/cobra"

func AddKernelCommand(root *cobra.Command) {
	root.AddCommand(initAppCommand())
	root.AddCommand(initCronCommand())
	root.AddCommand(initEnvCommand())
	root.AddCommand(initConfigCommand())
	root.AddCommand(initBuildCommand())
	root.AddCommand(initDevCommand())
	root.AddCommand(initProviderCommand())
	root.AddCommand(initCmdCommand())
	root.AddCommand(initMiddlewareCommand())
	root.AddCommand(initNewCommand())
	root.AddCommand(initSwaggerCommand())
}
