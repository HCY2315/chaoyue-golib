package log

import (
	"context"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
}

func (g *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	Infof(s, i...)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	Warnf(s, i...)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	Errorf(s, i...)
}

func (g *GormLogger) Trace(ctx context.Context, tm time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	Debugf("[rows:%d] %s", rows, sql)
}

func NewGormLogger() logger.Interface {
	return &GormLogger{}
}
