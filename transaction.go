package cql

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/logger"
)

// Transaction executes the function "toExec" inside a transaction
// The transaction is automatically rolled back in case "toExec" returns an error
// opts can be used to pass arguments to the transaction
func (db *DB) Transaction(
	ctx context.Context,
	toExec func(*DB) error,
	opts ...*sql.TxOptions,
) error {
	begin := time.Now()

	err := db.GormDB.Transaction(
		func(tx *gorm.DB) error {
			return toExec(&DB{
				GormDB:                tx,
				withLoggerFromContext: db.withLoggerFromContext,
			})
		},
		opts...,
	)
	if err != nil {
		return err
	}

	loggerInterface, isLoggerInterface := db.gormDBWithContext(ctx).Logger.(logger.Interface)
	if isLoggerInterface {
		loggerInterface.TraceTransaction(ctx, begin)
	}

	return nil
}
