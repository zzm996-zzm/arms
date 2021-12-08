package command

import "github.com/arms/framework/cobra"

func AddKernelCommand(root *cobra.Command) {
	root.AddCommand(initAppCommand())
}
