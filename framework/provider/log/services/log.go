package services

import (
	"context"
	"io"
	pkgLog "log"
	"time"

	"github.com/arms/framework"
	"github.com/arms/framework/contract"
	"github.com/arms/framework/provider/log/formatter"
)

type Log struct {
	level      contract.LogLevel
	formatter  contract.Formatter  // 日志格式化方法
	ctxFielder contract.CtxFielder // ctx获取上下文字段
	output     io.Writer           // 输出
	c          framework.Container // 容器
}

func (log *Log) IsLevelEnable(level contract.LogLevel) bool {
	return level <= log.level
}

func (log *Log) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	//先判断日志级别
	if !log.IsLevelEnable(level) {
		return nil
	}

	//使用ctxFielder 获取context中的信息
	fs := fields
	if log.ctxFielder(ctx) != nil {
		t := log.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}

	//如果绑定的trace服务，获取trace信息 TODO: trace

	//将日志信息按照formatter 序列化字符串
	if log.formatter == nil {
		log.formatter = formatter.TextFormatter
	}
	ct, err := log.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}

	// 如果是panic级别，则使用log进行panic
	if level == contract.PanicLevel {
		pkgLog.Panicln(string(ct))
		return nil
	}

	// 通过output进行输出
	log.output.Write(ct)
	log.output.Write([]byte("\r\n"))

	return nil
}

func (log *Log) SetOutPut(output io.Writer) {
	log.output = output
}

// Panic 输出panic的日志信息
func (log *Log) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.PanicLevel, ctx, msg, fields)
}

// Fatal will add fatal record which contains msg and fields
func (log *Log) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.FatalLevel, ctx, msg, fields)
}

// Error will add error record which contains msg and fields
func (log *Log) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.ErrorLevel, ctx, msg, fields)
}

// Warn will add warn record which contains msg and fields
func (log *Log) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.WarnLevel, ctx, msg, fields)
}

// Info 会打印出普通的日志信息
func (log *Log) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.InfoLevel, ctx, msg, fields)
}

// Debug will add debug record which contains msg and fields
func (log *Log) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.DebugLevel, ctx, msg, fields)
}

// Trace will add trace info which contains msg and fields
func (log *Log) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.TraceLevel, ctx, msg, fields)
}

// SetLevel set log level, and higher level will be recorded
func (log *Log) SetLevel(level contract.LogLevel) {
	log.level = level
}

// SetCxtFielder will get fields from context
func (log *Log) SetCtxFielder(handler contract.CtxFielder) {
	log.ctxFielder = handler
}

// SetFormatter will set formatter handler will covert data to string for recording
func (log *Log) SetFormatter(formatter contract.Formatter) {
	log.formatter = formatter
}

func (log *Log) SetContainer(c framework.Container) {
	log.c = c
}
