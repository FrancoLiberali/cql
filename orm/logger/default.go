package logger

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

const (
	defaultSlowQueryThreshold       = 200 * time.Millisecond
	defaultSlowTransactionThreshold = 200 * time.Millisecond
	DisableThreshold                = 0
)

var (
	DefaultConfig = Config{
		LogLevel:                  gormLogger.Warn,
		SlowQueryThreshold:        defaultSlowQueryThreshold,
		SlowTransactionThreshold:  defaultSlowTransactionThreshold,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
	}
	Default       = New(DefaultConfig)
	defaultWriter = log.New(os.Stdout, "\r\n", log.LstdFlags)
)

type defaultLogger struct {
	gormLogger.Interface
	Config
}

func New(config Config) Interface {
	return NewWithWriter(config, defaultWriter)
}

func NewWithWriter(config Config, writer gormLogger.Writer) Interface {
	return &defaultLogger{
		Config: config,
		Interface: gormLogger.New(
			writerWrapper{Writer: writer},
			config.toGormConfig(),
		),
	}
}

func (l *defaultLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	// method made to satisfy gormLogger.Interface
	return l.ToLogMode(level)
}

func (l *defaultLogger) ToLogMode(level gormLogger.LogLevel) Interface {
	newLogger := *l
	newLogger.LogLevel = level
	newLogger.Interface = newLogger.Interface.LogMode(level)

	return &newLogger
}

const nanoToMicro = 1e6

func (l defaultLogger) TraceTransaction(ctx context.Context, begin time.Time) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case l.SlowTransactionThreshold != DisableThreshold && elapsed > l.SlowTransactionThreshold && l.LogLevel >= gormLogger.Warn:
		l.Interface.Warn(ctx, "transaction_slow (>= %v) [%.3fms]", l.SlowTransactionThreshold, float64(elapsed.Nanoseconds())/nanoToMicro)
	case l.LogLevel >= gormLogger.Info:
		l.Interface.Info(ctx, "transaction_exec [%.3fms]", float64(elapsed.Nanoseconds())/nanoToMicro)
	}
}

type writerWrapper struct {
	Writer gormLogger.Writer
}

// Info, Warn, Error or Trace + Printf
const defaultStacktraceLen = 2

func (w writerWrapper) Printf(msg string, args ...interface{}) {
	if len(args) > 0 {
		// change the file path to avoid showing badaas-orm internal files
		firstArg := args[0]

		_, isString := firstArg.(string)
		if isString {
			file, line, caller := FindLastCaller(defaultStacktraceLen)
			if caller != 0 {
				w.Writer.Printf(
					msg,
					append(
						[]any{file + ":" + strconv.FormatInt(int64(line), 10)},
						args[1:]...,
					)...,
				)

				return
			}
		}
	}

	w.Writer.Printf(msg, args...)
}
