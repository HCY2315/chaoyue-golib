package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogConfig struct {
	hookFuncs []func(zapcore.Entry)
}

func newRotateFileCore(fileName string, maxDays, maxMB int, leveler zapcore.LevelEnabler) zapcore.Core {
	rotateWriter := &lumberjack.Logger{
		Filename: fileName,
		MaxSize:  maxMB,   // megabytes
		MaxAge:   maxDays, //days
	}
	fileEnc := newProductionEnc()
	return zapcore.NewCore(fileEnc, zapcore.AddSync(rotateWriter), leveler)
}

func newDevConsoleEnc() zapcore.Encoder {
	encConfig := zap.NewDevelopmentEncoderConfig()
	encConfig.EncodeTime = debugTimeEncoder
	encConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encConfig)
}

func debugTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.Kitchen))
}

func prodTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func newProductionEnc() zapcore.Encoder {
	enc := zap.NewProductionEncoderConfig()
	enc.EncodeTime = prodTimeEncoder
	return zapcore.NewConsoleEncoder(enc)
}

var levelToZapLevel = map[LogLevel]zapcore.Level{
	LogLevelDebug: zapcore.DebugLevel,
	LogLevelInfo:  zapcore.InfoLevel,
	LogLevelWarn:  zapcore.WarnLevel,
	LogLevelError: zapcore.ErrorLevel,
}

func buildLevelerFunc(level LogLevel) func() zapcore.LevelEnabler {
	zapLevel := zapcore.InfoLevel
	if l, find := levelToZapLevel[level]; find {
		zapLevel = l
	}
	return func() zapcore.LevelEnabler {
		return zap.NewAtomicLevelAt(zapLevel)
	}
}

func newDevConsoleCore() zapcore.Core {
	devEnc := newDevConsoleEnc()
	stdOutHook := zapcore.AddSync(os.Stdout)
	//stdErrHook := zapcore.AddSync(os.Stderr)
	return zapcore.NewCore(devEnc, stdOutHook, buildLevelerFunc(LogLevelDebug)())
}

func newProductionCores(fileName string, maxDays, maxMB int, level LogLevel) zapcore.Core {
	leveler := buildLevelerFunc(level)
	return newRotateFileCore(fileName, maxDays, maxMB, leveler())
}

func newDevCores(fileName string, maxDays, maxMB int) zapcore.Core {
	return zapcore.NewTee(
		newRotateFileCore(fileName, maxDays, maxMB, buildLevelerFunc(LogLevelDebug)()),
		newDevConsoleCore(),
	)
}

func newZapSugarLogger(cores ...zapcore.Core) *zap.SugaredLogger {
	logger := zap.New(zapcore.NewTee(cores...))
	return logger.Sugar()
}
