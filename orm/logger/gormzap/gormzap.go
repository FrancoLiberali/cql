package gormzap

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/ditrit/badaas/orm/logger"
)

// This type implement the [logger.Interface] interface.
// It is to be used as a replacement for the original logger
type gormzap struct {
	logger.Config
	ZapLogger *zap.Logger
}

// The constructor of the gormzap logger with default config
func NewDefault(zapLogger *zap.Logger) logger.Interface {
	return New(zapLogger, logger.DefaultConfig)
}

// The constructor of the gormzap logger
func New(zapLogger *zap.Logger, config logger.Config) logger.Interface {
	return &gormzap{
		ZapLogger: zapLogger,
		Config:    config,
	}
}

// Set the GORM's log mode to the value passed as argument
// Take into account that zap logger also have a log level
// that will determine if this logs are written or not
// GORM Info logs will generate a log with DebugLevel
// GORM Warn logs will generate a log with WarnLevel
// GORM Error logs will generate a log with ErrorLevel
func (l *gormzap) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	// method made to satisfy gormLogger.Interface
	return l.ToLogMode(level)
}

// Set the GORM's log mode to the value passed as argument
// Take into account that zap logger also have a log level
// that will determine if this logs are written or not
// GORM Info logs will generate a log with DebugLevel
// GORM Warn logs will generate a log with WarnLevel
// GORM Error logs will generate a log with ErrorLevel
func (l *gormzap) ToLogMode(level gormLogger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

// log info
func (l gormzap) Info(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		l.logger().Sugar().Debugf(str, args...)
	}
}

// log warning
func (l gormzap) Warn(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.logger().Sugar().Warnf(str, args...)
	}
}

// log an error
func (l gormzap) Error(_ context.Context, str string, args ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.logger().Sugar().Errorf(str, args...)
	}
}

// log a trace
func (l gormzap) Trace(
	_ context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsedTime := time.Since(begin)

	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rowsAffected := fc()
		l.logger().Error(
			"query_error",
			append(getZapFields(elapsedTime, rowsAffected, sql), zap.Error(err))...,
		)
	case l.SlowQueryThreshold != logger.DisableThreshold && elapsedTime > l.SlowQueryThreshold && l.LogLevel >= gormLogger.Warn:
		sql, rowsAffected := fc()
		l.logger().Warn(
			fmt.Sprintf("query_slow (>= %v)", l.SlowQueryThreshold),
			getZapFields(elapsedTime, rowsAffected, sql)...,
		)
	case l.LogLevel >= gormLogger.Info:
		sql, rowsAffected := fc()
		l.logger().Debug(
			"query_exec",
			getZapFields(elapsedTime, rowsAffected, sql)...,
		)
	}
}

func getZapFields(elapsedTime time.Duration, rowsAffected int64, sql string) []zapcore.Field {
	rowsAffectedString := strconv.FormatInt(rowsAffected, 10)
	if rowsAffected == -1 {
		rowsAffectedString = "-"
	}

	return []zapcore.Field{
		zap.Duration("elapsed_time", elapsedTime),
		zap.String("rows_affected", rowsAffectedString),
		zap.String("sql", sql),
	}
}

func (l gormzap) TraceTransaction(_ context.Context, begin time.Time) {
	elapsed := time.Since(begin)

	switch {
	case l.SlowTransactionThreshold != logger.DisableThreshold && elapsed > l.SlowTransactionThreshold && l.LogLevel >= gormLogger.Warn:
		l.logger().Warn(
			fmt.Sprintf("transaction_slow (>= %v)", l.SlowTransactionThreshold),
			zap.Duration("elapsed_time", elapsed),
		)
	case l.LogLevel >= gormLogger.Info:
		l.logger().Debug(
			"transaction_exec",
			zap.Duration("elapsed_time", elapsed),
		)
	}
}

// Filter parameters from queries depending of the value of ParameterizedQueries
func (l gormzap) ParamsFilter(_ context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.ParameterizedQueries {
		return sql, nil
	}

	return sql, params
}

// Info, Warn, Error or Trace + logger
const gormzapStacktraceLen = 2

// return a logger that log the right caller
func (l gormzap) logger() *zap.Logger {
	_, _, caller := logger.FindLastCaller(gormzapStacktraceLen)
	if caller == 0 {
		// in case we checked in all the stacktrace and none meet the conditions,
		// return the zap logger with the caller of gormzap, no matter where
		return l.ZapLogger.WithOptions(zap.AddCallerSkip(gormzapStacktraceLen))
	}

	return l.ZapLogger.WithOptions(zap.AddCallerSkip(caller - 1)) // -1 because here is how many we want to skip
}
