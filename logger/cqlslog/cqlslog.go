package cqlslog

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strconv"
	"time"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/FrancoLiberali/cql/logger"
)

// This type implement the [logger.Interface] interface.
// It is to be used as a replacement for the original logger
type cqlslog struct {
	logger.Config
	SLogger *slog.Logger
}

// The constructor of the cqlslog logger with default config
func NewDefault(sLogger *slog.Logger) logger.Interface {
	return New(sLogger, logger.DefaultConfig)
}

// The constructor of the cqlzap logger
func New(sLogger *slog.Logger, config logger.Config) logger.Interface {
	return &cqlslog{
		SLogger: sLogger,
		Config:  config,
	}
}

// Set the log mode to the value passed as argument
// Take into account that zap logger also have a log level
// that will determine if this logs are written or not
// Info logs will generate a log with DebugLevel
// Warn logs will generate a log with WarnLevel
// Error logs will generate a log with ErrorLevel
func (l *cqlslog) LogMode(level logger.LogLevel) gormLogger.Interface {
	// method made to satisfy gormLogger.Interface
	return l.ToLogMode(level)
}

// Set the GORM's log mode to the value passed as argument
// Take into account that zap logger also have a log level
// that will determine if this logs are written or not
// Info logs will generate a log with DebugLevel
// Warn logs will generate a log with WarnLevel
// Error logs will generate a log with ErrorLevel
func (l *cqlslog) ToLogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

// log info
func (l cqlslog) Info(ctx context.Context, str string, args ...any) {
	if l.LogLevel >= logger.Info {
		l.logAttrs(ctx, slog.LevelInfo, fmt.Sprintf(str, args...), nil)
	}
}

// log warning
func (l cqlslog) Warn(ctx context.Context, str string, args ...any) {
	if l.LogLevel >= logger.Warn {
		l.logAttrs(ctx, slog.LevelWarn, fmt.Sprintf(str, args...), nil)
	}
}

// log an error
func (l cqlslog) Error(ctx context.Context, str string, args ...any) {
	if l.LogLevel >= logger.Error {
		l.logAttrs(ctx, slog.LevelError, fmt.Sprintf(str, args...), nil)
	}
}

// log a trace
func (l cqlslog) Trace(
	ctx context.Context,
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
		l.logAttrs(
			ctx,
			slog.LevelError,
			"query_error",
			append(getSlogFields(elapsed, rowsAffected, sql), slog.Any("error", err))...,
		)
	case elapsed > l.SlowQueryThreshold && l.SlowQueryThreshold != logger.DisableThreshold && l.LogLevel >= logger.Warn:
		sql, rowsAffected := fc()
		l.logAttrs(
			ctx,
			slog.LevelWarn,
			fmt.Sprintf("query_slow (>= %v)", l.SlowQueryThreshold),
			getSlogFields(elapsed, rowsAffected, sql)...,
		)
	case l.LogLevel >= logger.Info:
		sql, rowsAffected := fc()
		l.logAttrs(
			ctx,
			slog.LevelDebug,
			"query_exec",
			getSlogFields(elapsed, rowsAffected, sql)...,
		)
	}
}

func getSlogFields(elapsedTime time.Duration, rowsAffected int64, sql string) []any {
	rowsAffectedString := strconv.FormatInt(rowsAffected, 10)
	if rowsAffected == -1 {
		rowsAffectedString = "-"
	}

	return []any{
		slog.Duration("elapsed_time", elapsedTime),
		slog.String("rows_affected", rowsAffectedString),
		slog.String("sql", sql),
	}
}

func (l cqlslog) TraceTransaction(ctx context.Context, begin time.Time) {
	elapsed := time.Since(begin)

	switch {
	case l.SlowTransactionThreshold != logger.DisableThreshold && elapsed > l.SlowTransactionThreshold && l.LogLevel >= logger.Warn:
		l.logAttrs(
			ctx,
			slog.LevelWarn,
			fmt.Sprintf("transaction_slow (>= %v)", l.SlowTransactionThreshold),
			slog.Duration("elapsed_time", elapsed),
		)
	case l.LogLevel >= logger.Info:
		l.logAttrs(
			ctx,
			slog.LevelDebug,
			"transaction_exec",
			slog.Duration("elapsed_time", elapsed),
		)
	}
}

// Filter parameters from queries depending of the value of ParameterizedQueries
func (l cqlslog) ParamsFilter(_ context.Context, sql string, params ...any) (string, []any) {
	if l.ParameterizedQueries {
		return sql, nil
	}

	return sql, params
}

// Info, Warn, Error or Trace + logAttrs
const cqlslogStacktraceLen = 2

// log adds context attributes and logs a message with the given slog level
func (l cqlslog) logAttrs(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	if !l.SLogger.Handler().Enabled(ctx, level) {
		return
	}

	_, _, caller := logger.FindLastCaller(cqlslogStacktraceLen)
	if caller == 0 {
		caller = cqlslogStacktraceLen
	}

	// Properly handle the PC for the caller
	var pc uintptr

	var pcs [1]uintptr

	runtime.Callers(caller, pcs[:])
	pc = pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(attrs...)

	_ = l.SLogger.Handler().Handle(ctx, r)
}
