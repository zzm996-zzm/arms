package cobra

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
)

type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

// SetContainer 设置服务容器
func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// GetContainer 获取容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

func (c *Command) SetParantNull() {
	c.parent = nil
}

//TODO:重复的方法剥离
func (c *Command) AddCronCommand(spec string, cmd *Command) {
	//crom结构是挂在在根Command上的
	root := c.Root()
	if root.Cron == nil {
		// 初始化Cron
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}
	//增加说明信息
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
	})

	//制作一个rootCommand
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParantNull()
	cronCmd.SetContainer(root.GetContainer())
	//增加调用函数
	root.Cron.AddFunc(spec, func() {
		//如果后续的command出现panic 这里要捕获
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			// 打印出err信息
			log.Println(err)
		}
	})
}

func (c *Command) AddDistributedCronCommand(serviceName string, spec string, cmd *Command, holdTime time.Duration) {
	root := c.Root()
	//初始化cron
	if root.Cron == nil {
		// 初始化Cron
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}

	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type:        "distributed-cron",
		Cmd:         cmd,
		Spec:        spec,
		ServiceName: serviceName,
	})

	appService := root.GetContainer().MustMake(contract.AppKey).(contract.ArmsApp)
	distributeServce := root.GetContainer().MustMake(contract.DistributedKey).(contract.Distributed)
	appID := appService.AppID()

	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParantNull()
	cronCmd.SetContainer(root.GetContainer())

	root.Cron.AddFunc(spec, func() {
		//防止panic
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		//节点选举
		selectAppID, err := distributeServce.Select(serviceName, appID, holdTime)
		if err != nil {
			return
		}

		//如果自己没有被选择到
		if selectAppID != appID {
			return
		}

		//如果被选择到了就执行
		err = cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})

}
