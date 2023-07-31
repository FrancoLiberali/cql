package gormdatabase

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/orm"
)

// Create the dsn string from the configuration
func createDialectorFromConf(databaseConfiguration configuration.DatabaseConfiguration) gorm.Dialector {
	return orm.CreateDialector(
		databaseConfiguration.GetHost(),
		databaseConfiguration.GetUsername(),
		databaseConfiguration.GetPassword(),
		databaseConfiguration.GetSSLMode(),
		databaseConfiguration.GetDBName(),
		databaseConfiguration.GetPort(),
	)
}

// Creates the database object with using the database configuration and exec the setup
func SetupDatabaseConnection(
	logger *zap.Logger,
	databaseConfiguration configuration.DatabaseConfiguration,
) (*gorm.DB, error) {
	dialector := createDialectorFromConf(databaseConfiguration)

	return orm.ConnectToDialector(
		logger,
		dialector,
		databaseConfiguration.GetRetry(),
		databaseConfiguration.GetRetryTime(),
	)
}
