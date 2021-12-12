package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/arms/framework"
	"github.com/arms/framework/cobra"
	"github.com/arms/framework/contract"
	"github.com/arms/framework/util"
	"github.com/pkg/errors"
)

// 初始化provider相关服务
func initProviderCommand() *cobra.Command {
	providerCommand.AddCommand(providerCreateCommand)
	providerCommand.AddCommand(providerListCommand)
	return providerCommand
}

var providerCommand = &cobra.Command{
	Use:   "provider",
	Short: "服务提供相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

var providerListCommand = &cobra.Command{
	Use:   "list",
	Short: "服务提供相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		app := container.(*framework.AppContainer)
		list := app.NameList()
		for _, name := range list {
			fmt.Println(name)
		}
		return nil
	},
}

var providerCreateCommand = &cobra.Command{
	Use:   "new",
	Short: "创建一个服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("创建一个服务")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入服务名称（服务凭证）:",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入服务所在目录名称(默认: 同服务名称):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}

			if folder == "" {
				folder = name
			}
		}

		//检查服务是否存在
		if container.IsBind(name) {
			fmt.Println("服务名称已经存在")
			return nil
		}
		app := container.MustMake(contract.AppKey).(contract.ArmsApp)
		pFolder := app.ProviderFolder()
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

		//创建文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}
		appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
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
		//创建title 这个模板方法
		funcs := template.FuncMap{"title": strings.Title}
		//创建 contract.go
		{
			file := filepath.Join(pFolder, folder, "contract.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("contract").Funcs(funcs).Parse(contractTmp))
			//将 name传递进入到template中渲染，并输出到contract.go
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}
		//创建provider.go
		{
			file := filepath.Join(pFolder, folder, "provider.go")
			f, err := os.Create(file)
			if err != nil {
				return nil
			}
			t := template.Must(template.New("provider").Funcs(funcs).Parse(providerTmp))
			if err := t.Execute(f, mu); err != nil {
				return err
			}
		}

		//创建service.go
		{
			file := filepath.Join(pFolder, folder, "service.go")
			f, err := os.Create(file)
			if err != nil {
				return err
			}
			t := template.Must(template.New("service").Funcs(funcs).Parse(serviceTmp))
			if err := t.Execute(f, mu); err != nil {
				return err
			}
		}
		fmt.Println("创建服务成功, 文件夹地址:", filepath.Join(pFolder, folder))
		fmt.Println("请不要忘记挂载新创建的服务")
		return nil
	},
}

var contractTmp string = `package {{.}}
const {{.|title}}Key = "{{.}}"
type Service interface {
	// 请在这里定义你的方法
    Foo() string
}
`

var providerTmp string = `package {{.Name}}
import (
	"{{.ModPath}}/framework"
)
type {{.Name|title}}Provider struct {
	framework.ServiceProvider
	c framework.Container
}
func (sp *{{.Name|title}}Provider) Name() string {
	return {{.Name|title}}Key
}
func (sp *{{.Name|title}}Provider) Register(c framework.Container) framework.NewInstance {
	return New{{.Name|title}}Service
}
func (sp *{{.Name|title}}Provider) IsDefer() bool {
	return false
}
func (sp *{{.Name|title}}Provider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}
func (sp *{{.Name|title}}Provider) Boot(c framework.Container) error {
	return nil
}
`
var serviceTmp string = `package {{.Name}}

import "{{.ModPath}}/framework"

type {{.Name|title}}Service struct {
	container framework.Container
}
func New{{.Name|title}}Service(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &{{.Name|title}}Service{container: container}, nil
}
func (s *{{.Name|title}}Service) Foo() string {
    return ""
}
`
