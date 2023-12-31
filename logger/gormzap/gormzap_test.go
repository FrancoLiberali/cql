package gormzap_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/FrancoLiberali/cql/logger"
	"github.com/FrancoLiberali/cql/logger/gormzap"
)

func TestTraceError(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(core)

	logger := gormzap.NewDefault(zapLogger)
	err := errors.New("sql error")
	logger.Trace(
		context.Background(),
		time.Now(),
		func() (string, int64) {
			return "fail sql", -1
		},
		err,
	)

	require.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, zapcore.ErrorLevel, log.Level)
	assert.Equal(t, "query_error", log.Message)
	require.Len(t, log.Context, 4)
	assert.Contains(t, log.Context, zap.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
	assert.Contains(t, log.Context, zap.Field{Key: "rows_affected", Type: zapcore.StringType, String: "-"})
	assert.Contains(t, log.Context, zap.Field{Key: "sql", Type: zapcore.StringType, String: "fail sql"})
}

func TestTraceSlowQuery(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(core)

	logger := gormzap.NewDefault(zapLogger)
	logger.Trace(
		context.Background(),
		time.Now().Add(-300*time.Millisecond),
		func() (string, int64) {
			return "slow sql", 1
		},
		nil,
	)

	require.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, zapcore.WarnLevel, log.Level)
	assert.Equal(t, "query_slow (>= 200ms)", log.Message)
	require.Len(t, log.Context, 3)
	assert.Contains(t, log.Context, zap.Field{Key: "rows_affected", Type: zapcore.StringType, String: "1"})
	assert.Contains(t, log.Context, zap.Field{Key: "sql", Type: zapcore.StringType, String: "slow sql"})
}

func TestTraceQueryExec(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(core)

	logger := gormzap.NewDefault(zapLogger).ToLogMode(logger.Info)
	logger.Trace(
		context.Background(),
		time.Now().Add(3*time.Hour),
		func() (string, int64) {
			return "normal sql", 1
		},
		nil,
	)

	require.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, zapcore.DebugLevel, log.Level)
	assert.Equal(t, "query_exec", log.Message)
	require.Len(t, log.Context, 3)
	assert.Contains(t, log.Context, zap.Field{Key: "rows_affected", Type: zapcore.StringType, String: "1"})
	assert.Contains(t, log.Context, zap.Field{Key: "sql", Type: zapcore.StringType, String: "normal sql"})
}

func TestTraceSlowTransaction(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(core)

	logger := gormzap.NewDefault(zapLogger)
	logger.TraceTransaction(
		context.Background(),
		time.Now().Add(-300*time.Millisecond),
	)

	require.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, zapcore.WarnLevel, log.Level)
	assert.Equal(t, "transaction_slow (>= 200ms)", log.Message)
	require.Len(t, log.Context, 1)
}

func TestTraceTransactionExec(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	zapLogger := zap.New(core)

	logger := gormzap.NewDefault(zapLogger).ToLogMode(logger.Info)
	logger.TraceTransaction(
		context.Background(),
		time.Now().Add(3*time.Hour),
	)

	require.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, zapcore.DebugLevel, log.Level)
	assert.Equal(t, "transaction_exec", log.Message)
	require.Len(t, log.Context, 1)
}
