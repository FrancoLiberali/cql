package database

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/orm"
	"github.com/ditrit/badaas/orm/logger"
	"github.com/ditrit/badaas/orm/logger/gormzap"
)

// Create the dsn string from the configuration
func createDialectorFromConf(databaseConfiguration configuration.DatabaseConfiguration) gorm.Dialector {
	return postgres.Open(orm.CreatePostgreSQLDSN(
		databaseConfiguration.GetHost(),
		databaseConfiguration.GetUsername(),
		databaseConfiguration.GetPassword(),
		databaseConfiguration.GetSSLMode(),
		databaseConfiguration.GetDBName(),
		databaseConfiguration.GetPort(),
	))
}

// Creates the database object with using the database configuration and exec the setup
func SetupDatabaseConnection(
	zapLogger *zap.Logger,
	databaseConfiguration configuration.DatabaseConfiguration,
	loggerConfiguration configuration.LoggerConfiguration,
) (*gorm.DB, error) {
	return OpenWithRetry(
		createDialectorFromConf(databaseConfiguration),
		gormzap.New(zapLogger, logger.Config{
			LogLevel:                  loggerConfiguration.GetLogLevel(),
			SlowQueryThreshold:        loggerConfiguration.GetSlowQueryThreshold(),
			SlowTransactionThreshold:  loggerConfiguration.GetSlowTransactionThreshold(),
			IgnoreRecordNotFoundError: loggerConfiguration.GetIgnoreRecordNotFoundError(),
			ParameterizedQueries:      loggerConfiguration.GetParameterizedQueries(),
		}),
		databaseConfiguration.GetRetry(),
		databaseConfiguration.GetRetryTime(),
	)
}

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
