package gormdatabase

import (
	"fmt"
	"time"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/gormdatabase/gormzap"
	"github.com/ditrit/badaas/persistence/models"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Create the dsn string from the configuration
func createDsnFromConf(databaseConfiguration configuration.DatabaseConfiguration) string {
	dsn := createDsn(
		databaseConfiguration.GetHost(),
		databaseConfiguration.GetUsername(),
		databaseConfiguration.GetPassword(),
		databaseConfiguration.GetSSLMode(),
		databaseConfiguration.GetDBName(),
		databaseConfiguration.GetPort(),
	)
	return dsn
}

// Create the dsn strings with the provided args
func createDsn(host, username, password, sslmode, dbname string, port int) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d sslmode=%s dbname=%s",
		username, password, host, port, sslmode, dbname,
	)
}

// Initialize the database with using the database configuration
func CreateDatabaseConnectionFromConfiguration(logger *zap.Logger, databaseConfiguration configuration.DatabaseConfiguration) (*gorm.DB, error) {
	dsn := createDsnFromConf(databaseConfiguration)
	var err error
	var database *gorm.DB
	for numberRetry := uint(0); numberRetry < databaseConfiguration.GetRetry(); numberRetry++ {
		database, err = initializeDBFromDsn(dsn, logger)
		if err == nil {
			logger.Sugar().Debugf("Database connection is active")
			err = AutoMigrate(logger, database)
			if err != nil {
				logger.Error("migration failed")
				return nil, err
			}
			logger.Info("AutoMigration was executed successfully")
			return database, err
		}
		logger.Sugar().Debugf("Database connection failed with error %q", err.Error())
		logger.Sugar().Debugf("Retrying database connection %d/%d in %s",
			numberRetry+1, databaseConfiguration.GetRetry(), databaseConfiguration.GetRetryTime().String())
		time.Sleep(databaseConfiguration.GetRetryTime())
	}
	return nil, err
}

// Initialize the database with the dsn string
func initializeDBFromDsn(dsn string, logger *zap.Logger) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormzap.New(logger),
	})

	if err != nil {
		return nil, err
	}

	rawDatabase, err := database.DB()
	if err != nil {
		return nil, err
	}
	// ping the underlying database
	err = rawDatabase.Ping()
	if err != nil {
		return nil, err
	}
	return database, nil
}

// Migrate the database using gorm [https://gorm.io/docs/migration.html#Auto-Migration]
func autoMigrate(database *gorm.DB, listOfDatabaseTables []any) error {
	err := database.AutoMigrate(listOfDatabaseTables...)
	if err != nil {
		return err
	}
	return nil
}

// Run the automigration
func AutoMigrate(logger *zap.Logger, database *gorm.DB) error {
	err := autoMigrate(database, models.ListOfTables)
	if err != nil {
		return err
	}
	return nil
}
