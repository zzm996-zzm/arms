package command

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/arms/framework/cobra"
	"github.com/arms/framework/contract"
)

// build相关的命令
func initBuildCommand() *cobra.Command {
	buildCommand.AddCommand(buildSelfCommand)
	buildCommand.AddCommand(buildBackendCommand)
	buildCommand.AddCommand(buildFrontendCommand)
	buildCommand.AddCommand(buildAllCommand)
	return buildCommand
}

var buildCommand = &cobra.Command{
	Use:   "build",
	Short: "编译相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var buildSelfCommand = &cobra.Command{
	Use:   "self",
	Short: "编译arms命令",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("arms go: please install go in path first")
		}

		cmd := exec.Command(path, "build", "-o", "arms", "-race", "./")
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("go build error:")
			fmt.Println(string(out))
			fmt.Println("--------------")
			return err
		}

		fmt.Println("build success please run ./arms")
		return nil
	},
}

var buildBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "使用go编译后端",
	RunE: func(c *cobra.Command, args []string) error {
		return buildSelfCommand.RunE(c, args)
	},
}

var buildFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "使用yarn编译前端",
	RunE: func(c *cobra.Command, args []string) error {
		path, err := exec.LookPath("yarn")
		appService := c.GetContainer().MustMake(contract.AppKey).(contract.ArmsApp)
		if err != nil {
			log.Fatalln("You need to install YARN under your path")
		}

		cmd := exec.Command(path, "run", "build")
		cmd.Dir = filepath.Join(appService.BaseFolder(), "vue")

		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("=============  前端编译失败 ============")
			fmt.Println(string(out))
			fmt.Println("=============  前端编译失败 ============")
			return err
		}

		// 打印输出
		fmt.Print(string(out))
		fmt.Println("=============  前端编译成功 ============")
		return nil
	},
}

var buildAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时编译前端和后端",
	RunE: func(c *cobra.Command, args []string) error {
		err := buildFrontendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		err = buildBackendCommand.RunE(c, args)
		if err != nil {
			return err
		}
		return nil
	},
}
