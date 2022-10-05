package gormdatabase

import (
	"fmt"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/persistence/gormdatabase/gormzap"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Create the dsn string from the configuration
func createDsnFromConf() string {
	databaseConfiguration := configuration.NewDatabaseConfiguration()
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

// Initialize the database with the configuration loaded by verdeter.
func InitializeDBFromConf() (*gorm.DB, error) {
	dsn := createDsnFromConf()
	return InitializeDBFromDsn(dsn)
}

// Initialize the database with the dsn string
func InitializeDBFromDsn(dsn string) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormzap.New(zap.L()),
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
func AutoMigrate(database *gorm.DB, listOfDatabaseTables ...any) error {
	err := database.AutoMigrate(listOfDatabaseTables...)
	if err != nil {
		return err
	}
	return nil
}
