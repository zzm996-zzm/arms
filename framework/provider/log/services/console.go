package services

import (
	"os"

	"github.com/zzm996-zzm/arms/framework"
	"github.com/zzm996-zzm/arms/framework/contract"
)

type ConsoleLog struct {
	Log
}

func NewConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &ConsoleLog{}
	log.SetFormatter(formatter)
	log.SetCtxFielder(ctxFielder)
	log.SetLevel(level)
	log.SetOutPut(os.Stdout)
	log.SetContainer(c)

	return log, nil
}
