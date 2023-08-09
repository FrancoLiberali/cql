package logger_test

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ditrit/badaas/orm/logger"
)

func TestTraceError(t *testing.T) {
	var buffer bytes.Buffer

	logger := logger.NewWithWriter(logger.DefaultConfig, log.New(&buffer, "\r\n", log.LstdFlags))
	err := errors.New("sql error")
	logger.Trace(
		context.Background(),
		time.Now(),
		func() (string, int64) {
			return "fail sql", -1
		},
		err,
	)

	assert.Contains(t, buffer.String(), "sql error")
	assert.Contains(t, buffer.String(), "[rows:-]")
	assert.Contains(t, buffer.String(), "fail sql")
}

func TestTraceSlowQuery(t *testing.T) {
	var buffer bytes.Buffer

	logger := logger.NewWithWriter(logger.DefaultConfig, log.New(&buffer, "\r\n", log.LstdFlags))
	logger.Trace(
		context.Background(),
		time.Now().Add(-300*time.Millisecond),
		func() (string, int64) {
			return "slow sql", 1
		},
		nil,
	)

	assert.Contains(t, buffer.String(), "SLOW SQL >= 200ms")
	assert.Contains(t, buffer.String(), "[rows:1]")
	assert.Contains(t, buffer.String(), "slow sql")
}

func TestTraceQueryExec(t *testing.T) {
	var buffer bytes.Buffer

	logger := logger.NewWithWriter(logger.DefaultConfig, log.New(&buffer, "\r\n", log.LstdFlags)).ToLogMode(logger.Info)
	logger.Trace(
		context.Background(),
		time.Now().Add(3*time.Hour),
		func() (string, int64) {
			return "normal sql", 1
		},
		nil,
	)

	assert.Contains(t, buffer.String(), "[rows:1]")
	assert.Contains(t, buffer.String(), "normal sql")
}
