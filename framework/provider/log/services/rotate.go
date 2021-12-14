package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
	"github.com/zzm996-zzm/arms/framework/util"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type RotateLog struct {
	Log
	//存储的目录
	folder string
	//文件名
	file string
}

func NewRotateLog(params ...interface{}) (interface{}, error) {
	// 参数解析
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	appService := c.MustMake(contract.AppKey).(contract.ArmsApp)
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	//获取folder
	folder := appService.LogFolder()
	if configService.IsExist("log.folder") {
		folder = configService.GetString("log.folder")
	}

	//如果folder 不存在
	//TODO:err handler
	if util.Exists(folder) {
		os.MkdirAll(folder, os.ModePerm)
	}

	// 从配置文件中获取file信息，否则使用默认的arms.log
	file := "arms.log"
	if configService.IsExist("log.file") {
		file = configService.GetString("log.file")
	}

	// 从配置文件获取date_format信息
	dateFormat := "%Y%m%d%H"
	if configService.IsExist("log.date_format") {
		dateFormat = configService.GetString("log.date_format")
	}

	linkName := rotatelogs.WithLinkName(filepath.Join(folder, file))
	options := []rotatelogs.Option{linkName}

	// 从配置文件获取rotate_count信息
	if configService.IsExist("log.rotate_count") {
		rotateCount := configService.GetInt("log.rotate_count")
		options = append(options, rotatelogs.WithRotationCount(uint(rotateCount)))
	}

	// 从配置文件获取rotate_size信息
	if configService.IsExist("log.rotate_size") {
		rotateSize := configService.GetInt("log.rotate_size")
		options = append(options, rotatelogs.WithRotationSize(int64(rotateSize)))
	}

	// 从配置文件获取max_age信息
	if configService.IsExist("log.max_age") {
		if maxAgeParse, err := time.ParseDuration(configService.GetString("log.max_age")); err == nil {
			options = append(options, rotatelogs.WithMaxAge(maxAgeParse))
		}
	}

	// 从配置文件获取rotate_time信息
	if configService.IsExist("log.rotate_time") {
		if rotateTimeParse, err := time.ParseDuration(configService.GetString("log.rotate_time")); err == nil {
			options = append(options, rotatelogs.WithRotationTime(rotateTimeParse))
		}
	}

	// 设置基础信息
	log := &RotateLog{}
	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)
	log.SetContainer(c)
	log.folder = folder
	log.file = file

	w, err := rotatelogs.New(fmt.Sprintf("%s.%s", filepath.Join(log.folder, log.file), dateFormat), options...)
	if err != nil {
		return nil, errors.Wrap(err, "new rotatelogs error")
	}
	log.SetOutput(w)
	return log, nil
}
