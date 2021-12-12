package command

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/go-github/v39/github"
	"github.com/spf13/cast"
	"github.com/zzm996-zzm/arms/framework/cobra"
	"github.com/zzm996-zzm/arms/framework/util"
)

func initNewCommand() *cobra.Command {
	return newCommand
}

var newCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个新应用",
	RunE: func(c *cobra.Command, args []string) error {
		currentPath := util.GetExecDirectory()
		var name string
		var folder string
		var mod string
		var version string
		var release *github.RepositoryRelease
		//目录名
		{
			prompt := &survey.Input{
				Message: "请输入目录名称：",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
			folder = path.Join(currentPath, name)
			if util.Exists(folder) {
				fmt.Println("目录" + folder + "已经存在")
				return nil
			}

		}
		//模块名
		{
			prompt := &survey.Input{
				Message: "请输入模块名称(go.mod中的module, 默认为文件夹名称)：",
			}
			err := survey.AskOne(prompt, &mod)
			if err != nil {
				return err
			}
			if mod == "" {
				mod = name
			}
		}
		//获取arms版本
		{
			//TODO:出现两次信息打印
			client := github.NewClient(nil)
			prompt := &survey.Input{
				Message: "请输入版本名称（参考 https://github.com/zzm996-zzm/arms/releases）默认最新版本",
			}
			err := survey.AskOne(prompt, &version)
			if err != nil {
				return err
			}
			if version != "" {
				//确认版本是否正确
				release, _, err = client.Repositories.GetReleaseByTag(context.Background(), "zzm996-zzm", "arms", version)
				if err != nil || release == nil {
					fmt.Println("版本不存在，创建失败，请参考 https://github.com/zzm996-zzm/arms/releases")
				}
			}
			if version == "" {
				release, _, err = client.Repositories.GetLatestRelease(context.Background(), "zzm996-zzm", "arms")
				version = release.GetTagName()
				if err != nil {
					fmt.Println("未知错误 创建失败，请参考 https://github.com/zzm996-zzm/arms/releases")
					return err
				}
			}
		}

		fmt.Println("====================================================")
		fmt.Println("开始进行创建应用操作")
		fmt.Println("创建目录：", folder)
		fmt.Println("应用名称：", mod)
		fmt.Println("arms框架版本：", release.GetTagName())

		templateName := "template-arms-" + version + "-" + cast.ToString(time.Now().Unix())
		templateFolder := path.Join(currentPath, templateName)
		os.Mkdir(templateFolder, os.ModePerm)
		fmt.Println("创建临时目录", templateFolder)

		//TODO:下载进度条
		//拷贝template项目
		url := release.GetZipballURL()
		err := util.DownloadFile(filepath.Join(templateFolder, "template.zip"), url)
		if err != nil {
			return err
		}

		fmt.Println("下载在zip包到template.zip")
		_, err = util.Unzip(filepath.Join(templateFolder, "template.zip"), templateFolder)
		if err != nil {
			return err
		}

		//获取folder下的arms 相关解压目录
		fInfos, err := ioutil.ReadDir(templateFolder)
		if err != nil {
			return err
		}

		for _, fInfo := range fInfos {
			//找到解压后的文件
			if fInfo.IsDir() && strings.Contains(fInfo.Name(), "arms") {
				//重命名
				if err := os.Rename(path.Join(templateFolder, fInfo.Name()), folder); err != nil {
					return err
				}
			}
		}

		fmt.Println("解压zip包")

		if err := os.RemoveAll(templateFolder); err != nil {
			return err
		}
		fmt.Println("删除临时文件夹", templateFolder)

		os.RemoveAll(filepath.Join(folder, ".git"))
		fmt.Println("删除.git目录")

		// 删除framework 目录
		os.RemoveAll(path.Join(folder, "framework"))
		fmt.Println("删除framework目录")

		filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			//TODO:go版本
			if path == filepath.Join(folder, "go.mod") {
				fmt.Println("更新文件:" + path)
				c = bytes.ReplaceAll(c, []byte("module github.com/zzm996-zzm/arms"), []byte("module "+mod))
				c = bytes.ReplaceAll(c, []byte("require ("), []byte("require (\n\tgithub.com/zzm996-zzm/arms "+version))
				err = ioutil.WriteFile(path, c, 0644)
				if err != nil {
					return err
				}
				return nil
			}

			isContain := bytes.Contains(c, []byte("github.com/zzm996-zzm/arms/app"))
			if isContain {
				fmt.Println("更新文件:" + path)
				c = bytes.ReplaceAll(c, []byte("github.com/zzm996-zzm/arms/app"), []byte(mod+"/app"))
				err = ioutil.WriteFile(path, c, 0644)
				if err != nil {
					return err
				}
			}

			return nil

		})
		fmt.Println("创建应用结束")
		fmt.Println("目录：", folder)
		fmt.Println("====================================================")
		return nil
	},
}
