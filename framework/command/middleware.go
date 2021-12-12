package command

import (
	"fmt"
	"io/ioutil"
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

// 初始化中间件相关命令
func initMiddlewareCommand() *cobra.Command {
	middlewareCommand.AddCommand(middlewareListCommand)
	// middlewareCommand.AddCommand(middlewareMigrateCommand)
	middlewareCommand.AddCommand(middlewareCreateCommand)
	return middlewareCommand
}

var middlewareCommand = &cobra.Command{
	Use:   "middleware",
	Short: "middleware相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

var middlewareListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有中间件",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
		middlewarePath := appService.MiddlewareFolder()
		//读取文件夹
		//TODO:重构
		files, err := ioutil.ReadDir(middlewarePath)
		if err != nil {
			return err
		}
		//仅仅打印文件夹名字，所有middleware由文件夹组成
		for _, f := range files {
			if f.IsDir() {
				fmt.Println(f.Name())
			}
		}
		return nil
	},
}

var middlewareCreateCommand = &cobra.Command{
	Use:   "new",
	Short: "创建middleware模板",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("开始创建middleware模板...")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入middleware名称命令",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认: 同middleware命令):",
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
		pFolder := appService.MiddlewareFolder()
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

		//创建模板
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
			t := template.Must(template.New("middleware").Funcs(funcs).Parse(middlewareTmp))
			// 将name传递进入到template中渲染，并且输出到contract.go 中
			if err := t.Execute(f, mu); err != nil {
				return errors.Cause(err)
			}
		}

		fmt.Println("创建middleware模板成功，路径:", filepath.Join(pFolder, folder))
		return nil
	},
}

var middlewareTmp string = `package {{.Name}}

import "{{.ModPath}}/framework/gin"

// {{.Name|title}}Middleware 代表中间件函数
func {{.Name|title}}Middleware() gin.HandlerFunc {

	return func(context *gin.Context) {
		context.Next()
	}

}
`
