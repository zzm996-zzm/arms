package command

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	cutil "github.com/zzm996-zzm/arms/framework/util/container"

	"github.com/erikdubbelboer/gspt"
	"github.com/sevlyar/go-daemon"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/cobra"
	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/util"
)

var appAddress string = ""
var appDaemon bool = false

const defaultAddress = ":8080"

func initAppCommand() *cobra.Command {
	appStartCommand.Flags().StringVar(&appAddress, "address", "", "设置app启动的地址，默认为:8080")
	appStartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", false, "start app daemon")
	appRestartCommand.Flags().BoolVarP(&appDaemon, "daemon", "d", true, "restart app daemon")
	appCommand.AddCommand(appStartCommand)
	appCommand.AddCommand(appStateCommand)
	appCommand.AddCommand(appStopCommand)
	appCommand.AddCommand(appRestartCommand)

	return appCommand
}

func getAppAddress(c framework.Container) string {
	if appAddress != "" {
		return appAddress
	}
	//分别从 env 和 config 中获取 优先级从前导后
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	if envService.Get("ADDRESS") != "" {
		appAddress = envService.Get("ADDRESS")
	} else if configService.IsExist("app.address") {
		appAddress = configService.GetString("app.address")
	} else {
		appAddress = defaultAddress
	}

	return appAddress
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

func getEngine(container framework.Container) http.Handler {
	// 从服务容器中获取kernel的服务实例
	kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
	// 从kernel服务实例中获取引擎
	return kernelService.HttpEngine()
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个Web服务",
	RunE: func(c *cobra.Command, args []string) error {
		var err error
		// 从Command中获取服务容器
		container := c.GetContainer()
		//获取http引擎
		core := getEngine(container)
		//获取端口
		address := getAppAddress(container)

		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    address,
		}

		serverPidFile, err := cutil.CreatePidFile(container, "app", os.Getpid())
		if err != nil {
			return err
		}
		serverLogFile, err := cutil.CreateLogFile(container, "app")
		if err != nil {
			return err
		}

		currentFolder := util.GetExecDirectory()
		//  daemon 模式
		if appDaemon {
			//创建一个Context
			cntxt := &daemon.Context{
				//设置pid文件
				PidFileName: serverPidFile,
				PidFilePerm: 0664,
				//日志文件
				LogFileName: serverLogFile,
				LogFilePerm: 0640,
				//设置工作路径
				WorkDir: currentFolder,
				//设置所有文件的mask，默认为750
				Umask: 027,
				//子进程参数
				Args: []string{"", "app", "start", "--daemon=true"},
			}
			// 启动子进程，d不为空表示当前是父进程，d为空表示当前是子进程
			d, err := cntxt.Reborn()
			if err != nil {
				return err
			}
			if d != nil {
				fmt.Println("启动成功，pid:", d.Pid)
				fmt.Println("日志文件:", serverLogFile)
				return nil
			}
			defer func(cntxt *daemon.Context) {
				err := cntxt.Release()
				if err != nil {
					fmt.Println(err.Error())
				}
			}(cntxt)
			//子进程执行真正的app启动操作
			fmt.Println("deamon started")
			gspt.SetProcTitle("arms app")
			if err := startAppServe(server, container); err != nil {
				fmt.Println(err)
			}
			return nil
		}

		// 非deamon模式，直接执行
		content := strconv.Itoa(os.Getpid())
		fmt.Println("[PID]", content)
		err = ioutil.WriteFile(serverPidFile, []byte(content), 0644)
		if err != nil {
			return err
		}
		gspt.SetProcTitle("arms app")

		fmt.Println("app serve url:", appAddress)
		if err := startAppServe(server, container); err != nil {
			fmt.Println(err)
		}
		return nil
	},
}

func startAppServe(server *http.Server, c framework.Container) error {
	// 这个goroutine是启动服务的goroutine
	go func() {
		server.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	return nil
}

func getPid(container framework.Container) (int, error) {
	appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
	pidFolder := appService.RuntimeFolder()
	pidFileName := filepath.Join(pidFolder, "app.pid")

	content, err := ioutil.ReadFile(pidFileName)
	if err != nil {
		return -1, err
	}
	if len(content) == 0 {
		fmt.Println("app service is stopd")
		return -1, errors.New("stoped")
	}
	//check
	pid, err := strconv.Atoi(string(content))
	if err != nil {
		return -1, errors.New("app pid file content no specification")
	}

	if !util.CheckProcessExist(pid) {
		return -1, errors.New("app service is stopd")
	}

	return pid, nil

}

var appStateCommand = &cobra.Command{
	Use:   "state",
	Short: "查看当前服务状态",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.ArmsApp)

		pidFolder := appService.RuntimeFolder()
		pidFileName := filepath.Join(pidFolder, "app.pid")

		content, err := ioutil.ReadFile(pidFileName)
		if err != nil || len(content) == 0 {
			return err
		}

		//check
		pid, err := strconv.Atoi(string(content))
		if err != nil {
			fmt.Println("app pid file content no specification")
			return err
		}
		if util.CheckProcessExist(pid) {
			fmt.Println("app service is running 【PID】:", pid)
			return nil
		}

		fmt.Println("app service not running ")
		return nil
	},
}

var appStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "停止当前服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
		pidFolder := appService.RuntimeFolder()
		pidFileName := filepath.Join(pidFolder, "app.pid")

		pid, err := getPid(container)
		fmt.Println("【PID】: ", pid)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			return err
		}

		if err := ioutil.WriteFile(pidFileName, []byte{}, 0644); err != nil {
			return err
		}

		fmt.Println("stop process 【PID】: ", pid)
		return nil
	},
}

var appRestartCommand = &cobra.Command{
	Use:   "restart",
	Short: "重启当前服务",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		pid, err := getPid(container)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		//关闭进程
		err = appStopCommand.RunE(c, args)
		if err != nil {
			return err
		}

		//判断是否真正的关闭了 15s轮询
		for i := 0; i < 15; i++ {
			if !util.CheckProcessExist(pid) {
				break
			}
			time.Sleep(time.Second * 1)
		}

		//再次检查是否退出
		if util.CheckProcessExist(pid) {
			return errors.New("stop error")
		}

		//重启服务
		fmt.Println("restart app service")
		err = appStartCommand.RunE(c, args)
		if err != nil {
			fmt.Println("start app service error")
			return err
		}

		return nil
	},
}
