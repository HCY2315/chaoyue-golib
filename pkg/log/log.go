// package log 提供日志功能
// Logger 是日志功能的抽象, 方便替换log实现
// 包级别方法Infof等，提供默认和可配置的Logger快捷方式

package log

import (
	"fmt"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/utils"
)

// ILogger 是日志的最小化功能, 以保证实现的可替代性
type ILogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

var defaultLogger ILogger

// init 初始化默认logger
// 只写入stdout而不是文件
// LogLevel设置为Debug
func init() {
	l := newDevelopmentLogger(defaultLogConfig)

	SetLogger(l)
}

// SetLogger 设置包级别的Logger代替defaultLogger
func SetLogger(logger ILogger) {
	defaultLogger = logger
	Debugf = defaultLogger.Debugf
	Infof = defaultLogger.Infof
	Warnf = defaultLogger.Warnf
	Errorf = defaultLogger.Errorf
}

type logFunc func(string, ...interface{})

var (
	Debugf logFunc
	Infof  logFunc
	Warnf  logFunc
	Errorf logFunc
)

type LogFileConfig struct {
	LogFileName  string
	MaxSizeInMB  int
	MaxAgeInDays int
	Level        string
}

var defaultLogFileConfig = LogFileConfig{
	LogFileName:  "./logs/default.log",
	MaxSizeInMB:  5,
	MaxAgeInDays: 3,
}

type LogConfig struct {
	RunMode    utils.RunMode
	Level      string
	FileConfig *LogFileConfig
}

var defaultLogConfig = LogConfig{
	RunMode:    utils.DevMode,
	FileConfig: &defaultLogFileConfig,
}

type LogOptions func(*LogConfig)

func WithLogRunMode(mode utils.RunMode) LogOptions {
	return func(lc *LogConfig) {
		lc.RunMode = mode
	}
}

func WithLogLevel(level string) LogOptions {
	return func(cfg *LogConfig) {
		cfg.Level = level
	}
}

func WithLogFile(lfc *LogFileConfig) LogOptions {
	if lfc == nil {
		lfc = &defaultLogFileConfig
	}
	return func(lc *LogConfig) {
		lc.FileConfig = lfc
	}
}

func NewLogger(options ...LogOptions) (ILogger, error) {
	lc := defaultLogConfig
	for _, op := range options {
		op(&lc)
	}
	if lc.RunMode == utils.DevMode {
		return newDevelopmentLogger(lc), nil
	}
	level, err := parseLogLevel(lc.Level)
	if err != nil {
		return nil, errors.Wrap(err, "解析logLevel:%s错误:%s", lc.Level)
	}
	return newProductionLogger(lc, level), nil
}

var logLevelSet = map[LogLevel]struct{}{
	LogLevelDebug: {},
	LogLevelInfo:  {},
	LogLevelWarn:  {},
	LogLevelError: {},
}

func parseLogLevel(level string) (LogLevel, error) {
	l := LogLevel(level)
	if _, find := logLevelSet[l]; find {
		return l, nil
	}
	return l, fmt.Errorf("未知LogLevel")
}

func newDevelopmentLogger(lc LogConfig) ILogger {
	devCores := newDevCores(lc.FileConfig.LogFileName, lc.FileConfig.MaxAgeInDays, lc.FileConfig.MaxSizeInMB)
	return newZapSugarLogger(devCores)
}

func newProductionLogger(lc LogConfig, logLevel LogLevel) ILogger {
	prdCores := newProductionCores(lc.FileConfig.LogFileName, lc.FileConfig.MaxAgeInDays, lc.FileConfig.MaxSizeInMB, logLevel)
	return newZapSugarLogger(prdCores)
}

type LogLevel string

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)
