package logger

import (
	"context"
	"io"
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
	// Discard logger will print any log to io.Discard
	Discard       = NewWithWriter(Config{}, log.New(io.Discard, "", log.LstdFlags))
	DefaultConfig = Config{
		LogLevel:                  Warn,
		SlowQueryThreshold:        defaultSlowQueryThreshold,
		SlowTransactionThreshold:  defaultSlowTransactionThreshold,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
	}
	// Default is default logger
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

func NewWithWriter(config Config, writer Writer) Interface {
	return &defaultLogger{
		Config: config,
		Interface: gormLogger.New(
			writerWrapper{Writer: writer},
			config.toGormConfig(),
		),
	}
}

func (l *defaultLogger) LogMode(level LogLevel) gormLogger.Interface {
	// method made to satisfy gormLogger.Interface
	return l.ToLogMode(level)
}

func (l *defaultLogger) ToLogMode(level LogLevel) Interface {
	newLogger := *l
	newLogger.LogLevel = level
	newLogger.Interface = newLogger.Interface.LogMode(level)

	return &newLogger
}

const nanoToMicro = 1e6

func (l defaultLogger) TraceTransaction(ctx context.Context, begin time.Time) {
	if l.LogLevel <= Silent {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case l.SlowTransactionThreshold != DisableThreshold && elapsed > l.SlowTransactionThreshold && l.LogLevel >= Warn:
		l.Interface.Warn(ctx, "transaction_slow (>= %v) [%.3fms]", l.SlowTransactionThreshold, float64(elapsed.Nanoseconds())/nanoToMicro)
	case l.LogLevel >= Info:
		l.Interface.Info(ctx, "transaction_exec [%.3fms]", float64(elapsed.Nanoseconds())/nanoToMicro)
	}
}

type writerWrapper struct {
	Writer Writer
}

// Info, Warn, Error or Trace + Printf
const defaultStacktraceLen = 2

func (w writerWrapper) Printf(msg string, args ...any) {
	if len(args) > 0 {
		// change the file path to avoid showing cql internal files
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
