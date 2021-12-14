package services

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/util"
)

type SingleLog struct {
	Log
	folder string
	file   string
	fd     *os.File
}

func NewSingleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	appService := c.MustMake(contract.AppKey).(contract.ArmsApp)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	log := &SingleLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetContainer(c)

	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}
	log.folder = folder
	if !util.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	log.file = "arms.log"
	if configService.IsExist("log.file") {
		log.file = configService.GetString("log.file")
	}

	fd, err := os.OpenFile(filepath.Join(log.folder, log.file), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "open log file err")
	}

	log.SetOutput(fd)
	log.c = c

	return log, nil
}
