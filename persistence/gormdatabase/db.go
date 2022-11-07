package gormdatabase

import (
	"fmt"

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
func InitializeDBFromConf(logger *zap.Logger, databaseConfiguration configuration.DatabaseConfiguration) (*gorm.DB, error) {
	dsn := createDsnFromConf(databaseConfiguration)
	db, err := initializeDBFromDsn(dsn, logger)
	if err != nil {
		return nil, err
	}
	err = autoMigrate(db, models.ListOfTables)
	if err != nil {
		return nil, err
	}
	logger.Info("The database connection was successfully initialized")
	return db, nil
}

// Initialize the database with the dsn string
func initializeDBFromDsn(dsn string, logger *zap.Logger) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormzap.New(logger),
	})

	if err != nil {
		return nil, err
	}

	rawDB, err := database.DB()
	if err != nil {
		return nil, err
	}
	// ping the underlying database
	err = rawDB.Ping()
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
