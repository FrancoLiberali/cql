package cqlzap

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

	"github.com/FrancoLiberali/cql/logger"
)

// This type implement the [logger.Interface] interface.
// It is to be used as a replacement for the original logger
type cqlzap struct {
	logger.Config
	ZapLogger *zap.Logger
}

// The constructor of the cqlzap logger with default config
func NewDefault(zapLogger *zap.Logger) logger.Interface {
	return New(zapLogger, logger.DefaultConfig)
}

// The constructor of the cqlzap logger
func New(zapLogger *zap.Logger, config logger.Config) logger.Interface {
	return &cqlzap{
		ZapLogger: zapLogger,
		Config:    config,
	}
}

// Set the log mode to the value passed as argument
// Take into account that zap logger also have a log level
// that will determine if this logs are written or not
// Info logs will generate a log with DebugLevel
// Warn logs will generate a log with WarnLevel
// Error logs will generate a log with ErrorLevel
func (l *cqlzap) LogMode(level logger.LogLevel) gormLogger.Interface {
	// method made to satisfy gormLogger.Interface
	return l.ToLogMode(level)
}

// Set the GORM's log mode to the value passed as argument
// Take into account that zap logger also have a log level
// that will determine if this logs are written or not
// Info logs will generate a log with DebugLevel
// Warn logs will generate a log with WarnLevel
// Error logs will generate a log with ErrorLevel
func (l *cqlzap) ToLogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

// log info
func (l cqlzap) Info(_ context.Context, str string, args ...any) {
	if l.LogLevel >= logger.Info {
		l.logger().Sugar().Debugf(str, args...)
	}
}

// log warning
func (l cqlzap) Warn(_ context.Context, str string, args ...any) {
	if l.LogLevel >= logger.Warn {
		l.logger().Sugar().Warnf(str, args...)
	}
}

// log an error
func (l cqlzap) Error(_ context.Context, str string, args ...any) {
	if l.LogLevel >= logger.Error {
		l.logger().Sugar().Errorf(str, args...)
	}
}

// log a trace
func (l cqlzap) Trace(
	_ context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rowsAffected := fc()
		l.logger().Error(
			"query_error",
			append(getZapFields(elapsed, rowsAffected, sql), zap.Error(err))...,
		)
	case elapsed > l.SlowQueryThreshold && l.SlowQueryThreshold != logger.DisableThreshold && l.LogLevel >= logger.Warn:
		sql, rowsAffected := fc()
		l.logger().Warn(
			fmt.Sprintf("query_slow (>= %v)", l.SlowQueryThreshold),
			getZapFields(elapsed, rowsAffected, sql)...,
		)
	case l.LogLevel >= logger.Info:
		sql, rowsAffected := fc()
		l.logger().Debug(
			"query_exec",
			getZapFields(elapsed, rowsAffected, sql)...,
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

func (l cqlzap) TraceTransaction(_ context.Context, begin time.Time) {
	elapsed := time.Since(begin)

	switch {
	case l.SlowTransactionThreshold != logger.DisableThreshold && elapsed > l.SlowTransactionThreshold && l.LogLevel >= logger.Warn:
		l.logger().Warn(
			fmt.Sprintf("transaction_slow (>= %v)", l.SlowTransactionThreshold),
			zap.Duration("elapsed_time", elapsed),
		)
	case l.LogLevel >= logger.Info:
		l.logger().Debug(
			"transaction_exec",
			zap.Duration("elapsed_time", elapsed),
		)
	}
}

// Filter parameters from queries depending of the value of ParameterizedQueries
func (l cqlzap) ParamsFilter(_ context.Context, sql string, params ...any) (string, []any) {
	if l.ParameterizedQueries {
		return sql, nil
	}

	return sql, params
}

// Info, Warn, Error or Trace + logger
const cqlzapStacktraceLen = 2

// return a logger that log the right caller
func (l cqlzap) logger() *zap.Logger {
	_, _, caller := logger.FindLastCaller(cqlzapStacktraceLen)
	if caller == 0 {
		// in case we checked in all the stacktrace and none meet the conditions,
		// return the zap logger with the caller of gormzap, no matter where
		return l.ZapLogger.WithOptions(zap.AddCallerSkip(cqlzapStacktraceLen))
	}

	return l.ZapLogger.WithOptions(zap.AddCallerSkip(caller - 1 - 1)) // -1 because here is how many we want to skip, -1 for runtime/proc.go:285
}
