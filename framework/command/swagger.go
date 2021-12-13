package command

import (
	"fmt"
	"github.com/swaggo/swag/gen"
	"github.com/zzm996-zzm/arms/framework/cobra"
	"github.com/zzm996-zzm/arms/framework/contract"
	"path"
)

func initSwaggerCommand() *cobra.Command{
	swaggerCommand.AddCommand(swaggerGenCommand)
	return swaggerCommand
}

var swaggerCommand = &cobra.Command{
	Use : "swagger",
	Short: "swagger相关命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

var swaggerGenCommand = &cobra.Command{
	Use:"gen",
	Short: "生成对应的swagger文件，contain swagger.yaml doc.go",
	Run: func(c *cobra.Command, args []string) {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
		httpFolder := appService.HttpFolder()
		outputDir := path.Join(httpFolder,"swagger")
		conf := &gen.Config{
			//遍历需要查询注释的目录
			SearchDir : httpFolder,
			//不包含哪些文件
			Excludes : "",
			//输出目录
			OutputDir: outputDir,
			//整个swagger接口说明文档注释
			MainAPIFile: "swagger.go",
			//名称显示的策略，比如首字母大写
			PropNamingStrategy: "",
			//是否要解析vendor目录
			ParseVendor: false,
			//是否要依赖外部库
			ParseDependency: false,
			//是否要解析标准包
			ParseInternal: false,
			//是否要查找markdown文件，这个markdown文件能用来为tag增加说明格式
			MarkdownFilesDir: "",
			//是否应该在docs.go中生成时间戳
			GeneratedTime: false,
		}

		err := gen.New().Build(conf)
		if err != nil {
			fmt.Println(err)
		}
	},
}