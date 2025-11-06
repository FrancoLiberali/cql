package cql

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"

	"github.com/FrancoLiberali/cql/logger"
	"github.com/FrancoLiberali/cql/logger/cqlslog"
	"github.com/FrancoLiberali/cql/test/models"
)

func TestWithLoggerFromContext(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer conn.Close()

	buffer := bytes.NewBuffer(nil)
	slogLogger := slog.New(slog.NewTextHandler(buffer, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	db, err := Open(postgres.New(postgres.Config{
		Conn: conn,
	}), WithLoggerFromContext(func(context.Context) logger.Interface {
		return cqlslog.NewDefault(slogLogger).ToLogMode(logger.Info)
	}))

	require.NoError(t, err)

	mock.ExpectBegin() // GORM might start a transaction
	mock.ExpectExec(`INSERT INTO "bicycles"`).
		WithArgs(sqlmock.AnyArg(), "John Doe", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit() // GORM might commit the transaction

	user := models.Bicycle{Name: "John Doe"}
	_, err = Insert(context.Background(), db, &user).Exec()
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())

	reader := bufio.NewReader(buffer)
	log, err := reader.ReadString('\n')
	require.NoError(t, err)

	assert.Contains(t, log, "time=")
	assert.Contains(t, log, "source=")
	assert.Contains(t, log, "cql/db_test.go:47")
	assert.Contains(t, log, "level=DEBUG")
	assert.Contains(t, log, `msg=query_exec`)
	assert.Contains(t, log, "elapsed_time=")
	assert.Contains(t, log, "rows_affected=1")
	assert.Contains(t, log, `sql="INSERT INTO \"bicycles\"`)
}
