package cutil

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/util"
)

func CreatePidFile(container framework.Container, filename string, pid int) (string, error) {

	appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
	pidFolder := appService.RuntimeFolder()

	//应用日志
	pidFilename := fmt.Sprintf("%s.pid", filename)
	serverPidFile := filepath.Join(pidFolder, pidFilename)
	return serverPidFile, util.CreateFile(pidFolder, pidFilename, []byte(strconv.Itoa(pid)))
}

func CreateLogFile(container framework.Container, filename string) (string, error) {
	appService := container.MustMake(contract.AppKey).(contract.ArmsApp)
	logFolder := appService.LogFolder()
	logFilename := fmt.Sprintf("%s.log", filename)
	serverLogFile := filepath.Join(logFolder, logFilename)
	return serverLogFile, util.CreateFile(logFolder, logFilename, []byte{})
}
