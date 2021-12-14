package orm

import (
	"context"
	"time"

	"github.com/zzm996-zzm/arms/framework/contract"
	"gorm.io/gorm/logger"
)

type OrmLogger struct {
	logger contract.Log
}

func NewOrmLogger(logger contract.Log) *OrmLogger {
	return &OrmLogger{logger: logger}
}

// Info 对接arms的info输出
func (o *OrmLogger) Info(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	o.logger.Info(ctx, s, fields)
}

// Warn 对接arms的Warn输出
func (o *OrmLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	o.logger.Warn(ctx, s, fields)
}

// Error 对接arms的Error输出
func (o *OrmLogger) Error(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	o.logger.Error(ctx, s, fields)
}

// Trace 对接arms的Trace输出
func (o *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	fields := map[string]interface{}{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}

	s := "orm trace sql"
	o.logger.Trace(ctx, s, fields)
}

// LogMode 什么都不实现，日志级别完全依赖hade的日志定义
func (o *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
	return o
}
