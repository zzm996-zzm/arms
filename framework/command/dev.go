package command

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/cobra"
	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/util"
)

func initDevCommand() *cobra.Command {
	devCommand.AddCommand(devBackendCommand)
	devCommand.AddCommand(devFrontendCommand)
	devCommand.AddCommand(devAllCommand)
	return devCommand
}

var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

var devBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "后端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		go proxy.monitorBackend()
		return proxy.startProxy(false, true)
	},
}

var devFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "前端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		return proxy.startProxy(true, false)
	},
}

var devAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时启动前端和后端调试",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		go proxy.monitorBackend()
		if err := proxy.startProxy(true, true); err != nil {
			return err
		}
		return nil
	},
}

type devConfig struct {
	//调试模式最终监听的端口
	Port string
	*Backend
	*Frontend
}

type Frontend struct {
	// 前端调试模式配置
	// 前端启动端口, 默认8071
	Port string
}

type Backend struct {
	// 调试模式后端更新时间，如果文件变更
	//等待3s才进行一次更新，能让频繁保存变更更为顺畅, 默认1s
	RefreshTime int
	//后端监听端口
	Port string
	//监听文件夹，默认为AppFolder
	MonitorFolder string
}

func newDeafultBackend() *Backend {
	return &Backend{
		1,
		"8072",
		"",
	}
}

func newDeafultFrontend() *Frontend {
	return &Frontend{
		"8071",
	}
}

func initDevConfig(c framework.Container) *devConfig {
	defaultPort := "8087"
	defaultBackend := newDeafultBackend()
	defaultFrontend := newDeafultFrontend()
	//设置默认值
	devConfig := &devConfig{
		Port:     defaultPort,
		Backend:  defaultBackend,
		Frontend: defaultFrontend,
	}

	configer := c.MustMake(contract.ConfigKey).(contract.Config)
	if configer.IsExist("app.dev.port") {
		devConfig.Port = configer.GetString("app.dev.port")
	}
	if configer.IsExist("app.dev.backend.refresh_time") {
		devConfig.Backend.RefreshTime = configer.GetInt("app.dev.backend.refresh_time")
	}
	if configer.IsExist("app.dev.backend.port") {
		devConfig.Backend.Port = configer.GetString("app.dev.backend.port")
	}
	monitorFolder := configer.GetString("app.dev.backend.monitor_folder")
	devConfig.Backend.MonitorFolder = monitorFolder
	if monitorFolder == "" {
		appService := c.MustMake(contract.AppKey).(contract.ArmsApp)
		devConfig.Backend.MonitorFolder = appService.AppFolder()
	}

	if configer.IsExist("app.dev.frontend.port") {
		devConfig.Frontend.Port = configer.GetString("app.dev.frontend.port")
	}
	return devConfig

}

type Proxy struct {
	devConfig   *devConfig
	backendPid  int //当前backend服务的pid
	frontendPid int //当前frontend服务的pid
	c           framework.Container
}

// NewProxy 初始化一个Proxy
func NewProxy(c framework.Container) *Proxy {
	devConfig := initDevConfig(c)
	return &Proxy{
		c:         c,
		devConfig: devConfig,
	}
}
func (p *Proxy) startProxy(startFrontend, startBackend bool) error {
	var backendURL, frontendURL *url.URL
	var err error
	if startFrontend {
		err = p.restartFrontend()
		if err != nil {
			return err
		}
	}

	if startBackend {
		err = p.restartBackend()
		if err != nil {
			return err
		}
	}

	frontendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Frontend.Port))
	if err != nil {
		return err
	}

	backendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Backend.Port))
	if err != nil {
		return err
	}

	proxyReverse := p.newProxyReverseProxy(frontendURL, backendURL)
	proxyServer := &http.Server{
		Addr:    "127.0.0.1:" + p.devConfig.Port,
		Handler: proxyReverse,
	}

	fmt.Println("代理服务启动:", "http://"+proxyServer.Addr)
	// 启动proxy服务
	err = proxyServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func (p *Proxy) newProxyReverseProxy(frontend, backend *url.URL) *httputil.ReverseProxy {
	if p.frontendPid == 0 && p.backendPid == 0 {
		fmt.Println("前端和后端服务都不存在")
		return nil
	}

	//仅仅存在后端
	if p.frontendPid == 0 {
		return httputil.NewSingleHostReverseProxy(backend)
	} else if p.backendPid == 0 {
		//仅仅存在前端
		return httputil.NewSingleHostReverseProxy(frontend)
	} else {
		//都存在
		director := func(req *http.Request) {
			if req.URL.Path == "/" || req.URL.Path == "/app.js" {
				req.URL.Scheme = frontend.Scheme
				req.URL.Host = frontend.Host
			} else {
				req.URL.Scheme = backend.Scheme
				req.URL.Host = backend.Host
			}
		}

		//定义一个 NotFoundErr
		NotFoundErr := errors.New("response is 404, need to redirect")
		return &httputil.ReverseProxy{
			Director: director,
			ModifyResponse: func(response *http.Response) error {
				//如果后端服务返回了404
				fmt.Println("状态码: ", response.StatusCode)
				if response.StatusCode == 404 {
					return NotFoundErr
				}
				return nil
			},
			ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
				if errors.Is(err, NotFoundErr) {
					httputil.NewSingleHostReverseProxy(frontend).ServeHTTP(writer, request)
				}
			},
		}
	}
}

func (p *Proxy) restartFrontend() error {
	if p.frontendPid != 0 {
		syscall.Kill(p.frontendPid, syscall.SIGKILL)
		p.frontendPid = 0
	}

	port := p.devConfig.Frontend.Port
	path, err := exec.LookPath("yarn")
	if err != nil {
		fmt.Println("you need install YARN to path")
		return err
	}
	cmd := exec.Command(path, "run", "serve")

	appService := p.c.MustMake(contract.AppKey).(contract.ArmsApp)
	cmd.Dir = filepath.Join(appService.BaseFolder(), "vue")

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s%s", "PORT=", port))
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	fmt.Println("启动前端服务: ", "http://127.0.0.1:"+port)
	if err != nil {
		//TODO:不处理错误吗？
		fmt.Println(err)
	}
	p.frontendPid = cmd.Process.Pid
	fmt.Println("前端服务pid:", p.frontendPid)
	return nil
}
func (p *Proxy) restartBackend() error {
	//杀死之前的进程
	if p.backendPid != 0 {
		syscall.Kill(p.backendPid, syscall.SIGKILL)
		p.backendPid = 0
	}

	//设置端口
	port := p.devConfig.Backend.Port
	address := fmt.Sprintf(":" + port)
	cmd := exec.Command("./arms", "app", "start", "--address="+address)
	cmd.Stdout = os.Stdin
	cmd.Stderr = os.Stderr
	fmt.Println("启动后端服务: ", "http://127.0.0.1"+address)

	err := cmd.Start()
	if err != nil {
		//TODO:no handle err?
		fmt.Println(err)
	}

	p.backendPid = cmd.Process.Pid
	fmt.Println("后端服务pid: ", p.backendPid)
	return nil
}

func (p *Proxy) rebuildBackend() error {
	//重新编译
	cmdBuild := exec.Command("./arms", "build", "backend")
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if err := cmdBuild.Start(); err == nil {
		err = cmdBuild.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Proxy) monitorBackend() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcher.Close()

	//开启监听目标文件夹
	appFolder := p.devConfig.Backend.MonitorFolder
	fmt.Println("监控文件夹：", appFolder)
	filepath.Walk(appFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}
		// 如果是隐藏的目录比如 . 或者 .. 则不用进行监控
		if util.IsHiddenDirectory(path) {
			return nil
		}

		return watcher.Add(path)
	})

	//开启计时机制
	refreshTime := p.devConfig.Backend.RefreshTime
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	//先停止计时器
	t.Stop()
	for {
		select {
		case <-t.C:
			// 计时器时间到了，代表之前有文件更新事件重置过计时器
			fmt.Println("检测文件更新，重启服务开始....")
			if err := p.rebuildBackend(); err != nil {
				fmt.Println("重新编译失败: ", err.Error())
			} else {
				if err := p.restartBackend(); err != nil {
					fmt.Println("重新启动失败：", err.Error())
				}
			}
			fmt.Println("重启服务结束....")
			t.Stop()
		case event, ok := <-watcher.Events:
			fmt.Println("event:", event.Name, event.Op.String())
			if !ok {
				continue
			}
			if strings.Contains(event.Op.String(),"WRITE"){
				//如果有文件更新reset 计时器
				fmt.Println("event:", event.Name, event.Op.String())
				t.Reset(time.Duration(refreshTime) * time.Second)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}

			fmt.Println("监听文件夹错误：", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}
	}
}
