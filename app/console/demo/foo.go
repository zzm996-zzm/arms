package demo

import (
	"log"

	"github.com/arms/framework/cobra"
)

func InitFoo() *cobra.Command {
	FooCommand.AddCommand(Foo1Command)
	return FooCommand
}

var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo的简要说明",
	Long:    "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		log.Println(container)
		log.Println("im foo")
		return nil
	},
}

// Foo1Command 代表Foo命令的子命令Foo1
var Foo1Command = &cobra.Command{
	Use:     "foo1",
	Short:   "foo1的简要说明",
	Long:    "foo1的长说明",
	Aliases: []string{"fo1", "f1"},
	Example: "foo1命令的例子",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		log.Println(container)
		log.Println("im foo1")
		return nil
	},
}
