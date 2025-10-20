package cqlslog_test

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/FrancoLiberali/cql/logger"
	"github.com/FrancoLiberali/cql/logger/cqlslog"
)

func TestTraceError(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	slogLogger := slog.New(slog.NewTextHandler(buffer, &slog.HandlerOptions{AddSource: true}))

	logger := cqlslog.NewDefault(slogLogger)
	err := errors.New("sql error")
	logger.Trace(
		context.Background(),
		time.Now(),
		func() (string, int64) {
			return "fail sql", -1
		},
		err,
	)

	reader := bufio.NewReader(buffer)
	log, err := reader.ReadString('\n')
	require.NoError(t, err)

	assert.Contains(t, log, "time=")
	assert.Contains(t, log, "source=")
	assert.Contains(t, log, "cql/logger/cqlslog/cqlslog_test.go:25")
	assert.Contains(t, log, "level=ERROR")
	assert.Contains(t, log, "msg=query_error")
	assert.Contains(t, log, "elapsed_time=")
	assert.Contains(t, log, "rows_affected=-")
	assert.Contains(t, log, `sql="fail sql"`)
	assert.Contains(t, log, `error="sql error"`)
}

func TestTraceSlowQuery(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	slogLogger := slog.New(slog.NewTextHandler(buffer, &slog.HandlerOptions{AddSource: true}))

	logger := cqlslog.NewDefault(slogLogger)
	logger.Trace(
		context.Background(),
		time.Now().Add(-300*time.Millisecond),
		func() (string, int64) {
			return "slow sql", 1
		},
		nil,
	)

	reader := bufio.NewReader(buffer)
	log, err := reader.ReadString('\n')
	require.NoError(t, err)

	assert.Contains(t, log, "time=")
	assert.Contains(t, log, "source=")
	assert.Contains(t, log, "cql/logger/cqlslog/cqlslog_test.go:54")
	assert.Contains(t, log, "level=WARN")
	assert.Contains(t, log, `msg="query_slow (>= 200ms)"`)
	assert.Contains(t, log, "elapsed_time=")
	assert.Contains(t, log, "rows_affected=1")
	assert.Contains(t, log, `sql="slow sql"`)
}

func TestTraceQueryExec(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	slogLogger := slog.New(slog.NewTextHandler(buffer, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	logger := cqlslog.NewDefault(slogLogger).ToLogMode(logger.Info)
	logger.Trace(
		context.Background(),
		time.Now().Add(3*time.Hour),
		func() (string, int64) {
			return "normal sql", 1
		},
		nil,
	)

	reader := bufio.NewReader(buffer)
	log, err := reader.ReadString('\n')
	require.NoError(t, err)

	assert.Contains(t, log, "time=")
	assert.Contains(t, log, "source=")
	assert.Contains(t, log, "cql/logger/cqlslog/cqlslog_test.go:85")
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, `msg=query_exec`)
	assert.Contains(t, log, "elapsed_time=")
	assert.Contains(t, log, "rows_affected=1")
	assert.Contains(t, log, `sql="normal sql"`)
}

func TestTraceSlowTransaction(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	slogLogger := slog.New(slog.NewTextHandler(buffer, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	logger := cqlslog.NewDefault(slogLogger)
	logger.TraceTransaction(
		context.Background(),
		time.Now().Add(-300*time.Millisecond),
	)

	reader := bufio.NewReader(buffer)
	log, err := reader.ReadString('\n')
	require.NoError(t, err)

	assert.Contains(t, log, "time=")
	assert.Contains(t, log, "source=")
	assert.Contains(t, log, "cql/logger/cqlslog/cqlslog_test.go:116")
	assert.Contains(t, log, "level=WARN")
	assert.Contains(t, log, `msg="transaction_slow (>= 200ms)"`)
	assert.Contains(t, log, "elapsed_time=")
}

func TestTraceTransactionExec(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	slogLogger := slog.New(slog.NewTextHandler(buffer, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	logger := cqlslog.NewDefault(slogLogger).ToLogMode(logger.Info)
	logger.TraceTransaction(
		context.Background(),
		time.Now().Add(3*time.Hour),
	)

	reader := bufio.NewReader(buffer)
	log, err := reader.ReadString('\n')
	require.NoError(t, err)

	assert.Contains(t, log, "time=")
	assert.Contains(t, log, "source=")
	assert.Contains(t, log, "cql/logger/cqlslog/cqlslog_test.go:141")
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, `msg=transaction_exec`)
	assert.Contains(t, log, "elapsed_time=")
}
