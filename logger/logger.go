package logger

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type Interface interface {
	gormLogger.Interface
	// change log mode
	ToLogMode(level gormLogger.LogLevel) Interface
	// Trace a committed transaction
	TraceTransaction(ctx context.Context, begin time.Time)
}

const (
	// Silent silent log level
	Silent gormLogger.LogLevel = gormLogger.Silent
	// Error error log level
	Error gormLogger.LogLevel = gormLogger.Error
	// Warn warn log level
	Warn gormLogger.LogLevel = gormLogger.Warn
	// Info info log level
	Info gormLogger.LogLevel = gormLogger.Info
)

type Config struct {
	LogLevel                  gormLogger.LogLevel // GORM's Log level: the level of the logs generated by gorm
	SlowQueryThreshold        time.Duration       // Slow SQL Query threshold (use DisableThreshold to disable it)
	SlowTransactionThreshold  time.Duration       // Slow Transaction threshold (use DisableThreshold to disable it)
	IgnoreRecordNotFoundError bool                // if true, ignore gorm.ErrRecordNotFound error for logger
	ParameterizedQueries      bool                // if true, don't include params in the query execution logs
	Colorful                  bool                // log with colors
}

func (c Config) toGormConfig() gormLogger.Config {
	return gormLogger.Config{
		LogLevel:                  c.LogLevel,
		SlowThreshold:             c.SlowQueryThreshold,
		IgnoreRecordNotFoundError: c.IgnoreRecordNotFoundError,
		ParameterizedQueries:      c.ParameterizedQueries,
		Colorful:                  c.Colorful,
	}
}

// search in the stacktrace the last file outside gormzap, cql and gorm
func FindLastCaller(skip int) (string, int, int) {
	// +1 because at least one will be inside gorm
	// +1 because of this function
	for i := skip + 1 + 1; i < 18; i++ {
		_, file, line, ok := runtime.Caller(i)

		if !ok {
			// we checked in all the stacktrace and none meet the conditions,
			return "", 0, 0
		} else if !strings.Contains(file, gormSourceDir) && !strings.Contains(file, gormForkSourceDir) && !strings.Contains(file, cqlSourceDir) {
			// file outside cql and gorm
			return file, line, i - 1 // -1 to remove this function from the stacktrace
		}
	}

	return "", 0, 0
}

var (
	cqlSourceDir      string
	gormSourceDir     = filepath.Join("gorm.io", "gorm")
	gormForkSourceDir = filepath.Join("!franco!liberali", "gorm")
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get cql source directory with various operating systems
	cqlSourceDir = sourceDir(file)
}

func sourceDir(file string) string {
	loggerDir := filepath.Dir(file)
	cqlDir := filepath.Dir(loggerDir)

	return filepath.ToSlash(cqlDir) + "/"
}
