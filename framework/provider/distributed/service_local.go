package distributed

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
)

type LocalDistributedService struct {
	container framework.Container
}

func NewLocalDistributedService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &LocalDistributedService{container: container}, nil
}

func (local *LocalDistributedService) Select(serviceName string, appID string, holdTime time.Duration) (selectAppID string, err error) {
	appService := local.container.MustMake(contract.AppKey).(contract.ArmsApp)
	runtimeFolder := appService.RuntimeFolder()
	localFile := filepath.Join(runtimeFolder, "distribute_"+serviceName)

	//打开文件锁
	lock, err := os.OpenFile(localFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	//尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	//抢不到文件
	if err != nil {
		//读取被选择的appid
		selectAppIDByt, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByt), err
	}

	//在一段时间内选举有效
	go func() {
		defer func() {
			//释放文件锁
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			//释放文件
			lock.Close()
			//删除文件锁对应的文件
			os.Remove(localFile)

		}()

		timer := time.NewTimer(holdTime)
		//等待计时器结束
		<-timer.C
	}()

	//抢占到了
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}

	return appID, nil
}
