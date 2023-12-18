package orm

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/logger"
)

// Executes the function "toExec" inside a transaction
// The transaction is automatically rolled back in case "toExec" returns an error
// opts can be used to pass arguments to the transaction
func Transaction[RT any](
	db *gorm.DB,
	toExec func(*gorm.DB) (RT, error),
	opts ...*sql.TxOptions,
) (RT, error) {
	begin := time.Now()

	var returnValue RT

	err := db.Transaction(
		func(tx *gorm.DB) error {
			var err error
			returnValue, err = toExec(tx)

			return err
		},
		opts...,
	)
	if err != nil {
		return *new(RT), err
	}

	loggerInterface, isLoggerInterface := db.Logger.(logger.Interface)
	if isLoggerInterface {
		loggerInterface.TraceTransaction(context.Background(), begin)
	}

	return returnValue, nil
}
