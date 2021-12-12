package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/arms/framework/cobra"
	"github.com/arms/framework/contract"
	"github.com/arms/framework/util"
	"github.com/pkg/errors"
)

// 初始化command相关命令
func initCmdCommand() *cobra.Command {
	_Command.AddCommand(listCommand)
	_Command.AddCommand(createCommand)
	return _Command
}

var _Command = &cobra.Command{
	Use:   "command",
	Short: "系统command相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

// listCommand 列出所有的控制台命令
var listCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有控制台命令",
	RunE: func(c *cobra.Command, args []string) error {
		cmds := c.Root().Commands()
		ps := [][]string{}
		for _, cmd := range cmds {
			line := []string{cmd.Name(), cmd.Short}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

var createCommand = &cobra.Command{
	Use:   "new",
	Short: "创建一个控制台命令",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("开始创建控制命令...")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入控制台命令",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认: 同控制台命令):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}
		if folder == "" {
			folder = name
		}

		appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
		pFolder := appService.CommandFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		for i := 0; i < len(subFolders); i++ {
			if subFolders[i] == folder {
				fmt.Println("目录名称已经存在")
				return nil
			}
		}

		modPath := util.GetModule(filepath.Join(appService.BaseFolder(), "go.mod"))
		modPath = strings.ReplaceAll(modPath, "module", "")
		modPath = strings.TrimSpace(modPath)
		//模板需要的字段
		mu := struct {
			Name    string
			ModPath string
		}{
			Name:    name,
			ModPath: modPath,
		}

		// 开始创建文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}

		// 创建title这个模版方法
		funcs := template.FuncMap{"title": strings.Title}
		{
			//  创建name.go
			file := filepath.Join(pFolder, folder, name+".go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			// 使用contractTmp模版来初始化template，并且让这个模版支持title方法，即支持{{.|title}}
			t := template.Must(template.New("cmd").Funcs(funcs).Parse(cmdTmpl))
			// 将name传递进入到template中渲染，并且输出到contract.go 中
			if err := t.Execute(f, mu); err != nil {
				return errors.Cause(err)
			}
		}

		fmt.Println("创建新命令行工具成功，路径:", filepath.Join(pFolder, folder))
		fmt.Println("请记得开发完成后将命令行工具挂载到 console/kernel.go")
		return nil

	},
}

// 命令行工具模版
var cmdTmpl string = `package {{.Name}}

import (
   "fmt"

   "{{.ModPath}}/framework/cobra"
)

var {{.Name}}Command = &cobra.Command{
   Use:   "{{.Name}}",
   Short: "{{.Name}}",
   RunE: func(c *cobra.Command, args []string) error {
      container := c.GetContainer()
      fmt.Println(container)
      return nil
   },
}`
