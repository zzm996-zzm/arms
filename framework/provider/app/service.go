package app

import (
	"errors"
	"path/filepath"

	"github.com/arms/framework"
	"github.com/arms/framework/contract"
	"github.com/arms/framework/util"
	"github.com/google/uuid"
)

var _ contract.ArmsApp

type ArmsApp struct {
	contract.ArmsApp
	container  framework.Container //服务容器
	appId      string
	baseFolder string //基础路径
}

func (app *ArmsApp) AppID() string {
	return app.appId
}

func (app *ArmsApp) Version() string {
	return "0.0.1"
}

func (app *ArmsApp) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}

	// //如果没有设置
	// var baseFolder string
	// flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
	// flag.Parse()
	// if baseFolder != "" {
	// 	return baseFolder
	// }

	//如果参数也没有，则使用默认的当前路径
	return util.GetExecDirectory()
}

func (app *ArmsApp) ConfigFolder() string {
	return filepath.Join(app.BaseFolder(), "config")
}

func (app ArmsApp) StorageFolder() string {
	return filepath.Join(app.BaseFolder(), "storage")
}
func (app ArmsApp) HttpFolder() string {
	return filepath.Join(app.BaseFolder(), "http")
}

func (app ArmsApp) ConsoleFolder() string {
	return filepath.Join(app.BaseFolder(), "console")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (app ArmsApp) ProviderFolder() string {
	return filepath.Join(app.BaseFolder(), "provider")
}

// TestFolder 定义测试需要的信息
func (app ArmsApp) TestFolder() string {
	return filepath.Join(app.BaseFolder(), "test")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (app ArmsApp) MiddlewareFolder() string {
	return filepath.Join(app.HttpFolder(), "middleware")
}

// LogFolder 表示日志存放地址
func (app ArmsApp) LogFolder() string {
	return filepath.Join(app.StorageFolder(), "log")
}

// CommandFolder 定义业务定义的命令
func (app ArmsApp) CommandFolder() string {
	return filepath.Join(app.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (app ArmsApp) RuntimeFolder() string {
	return filepath.Join(app.StorageFolder(), "runtime")
}

// NewHadeApp 初始化HadeApp
func NewArmsApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	appId := uuid.New().String()
	return &ArmsApp{appId: appId, baseFolder: baseFolder, container: container}, nil
}
