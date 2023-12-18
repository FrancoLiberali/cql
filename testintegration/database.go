package testintegration

import (
	"context"
	"time"

	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/logger"
	"gorm.io/gorm"
)

func OpenWithRetry(
	dialector gorm.Dialector,
	logger logger.Interface,
	connectionTries uint,
	retryTime time.Duration,
) (*gorm.DB, error) {
	var err error

	var gormDB *gorm.DB

	for retryNumber := uint(0); retryNumber < connectionTries; retryNumber++ {
		gormDB, err = orm.Open(
			dialector,
			&gorm.Config{
				Logger: logger,
			},
		)

		if err == nil {
			logger.Info(context.Background(), "Database connection is active")

			return gormDB, nil
		}

		// there are more retries
		if retryNumber < connectionTries-1 {
			logger.Info(
				context.Background(),
				"Database connection failed with error %q, retrying %d/%d in %s",
				err.Error(),
				retryNumber+1+1, // +1 for counting from 1 and +1 for next iteration
				connectionTries,
				retryTime,
			)
			time.Sleep(retryTime)
		} else {
			logger.Error(
				context.Background(),
				"Database connection failed with error %q",
				err.Error(),
			)
		}
	}

	return nil, err
}
