package db

import (
	"fmt"

	"github.com/ditrit/badaas/configuration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// the list of models (eg tables for the database) to use in migrations.
var listOfDatabaseTables = []any{}

// Database instance for the whole app
var localDB *gorm.DB

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
func InitializeDBFromConf() error {
	dsn := createDsnFromConf()
	return InitializeDBFromDsn(dsn)
}

// Initialize the database with the dsn string
func InitializeDBFromDsn(dsn string) error {
	var err error
	localDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("can't create the gorm.DB (%w)", err)
	}

	rawDB, err := localDB.DB()
	if err != nil {
		return fmt.Errorf("the database was not initialized correctly (%w)", err)
	}
	// ping the underlying database
	err = rawDB.Ping()
	if err != nil {
		return fmt.Errorf("can't ping the database (%w)", err)
	}
	return nil
}

// Return the database.
//
// Panics if the database is nil.
func GetDB() *gorm.DB {
	if localDB == nil {
		panic("the Database was de-allocated")
	}
	return localDB
}

func AutoMigrate() error {
	db := GetDB()
	return db.AutoMigrate(listOfDatabaseTables...)
}
