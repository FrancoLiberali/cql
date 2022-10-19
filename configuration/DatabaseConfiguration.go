package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Hold the configuration values for the database connection
type DatabaseConfiguration interface {
	ConfigurationHolder
	GetPort() int
	GetHost() string
	GetDBName() string
	GetUsername() string
	GetPassword() string
	GetSSLMode() string
}

// Concrete implementation of the DatabaseConfiguration interface
type databaseConfigurationImpl struct {
	port     int
	host     string
	dbName   string
	username string
	password string
	sslmode  string
}

// Instantiate a new configuration holder for the database connection
func NewDatabaseConfiguration() DatabaseConfiguration {
	databaseConfiguration := new(databaseConfigurationImpl)
	databaseConfiguration.Reload()
	return databaseConfiguration
}

// Reload database configuration
func (databaseConfiguration *databaseConfigurationImpl) Reload() {
	databaseConfiguration.port = viper.GetInt("database.port")
	databaseConfiguration.host = viper.GetString("database.host")
	databaseConfiguration.dbName = viper.GetString("database.name")
	databaseConfiguration.username = viper.GetString("database.username")
	databaseConfiguration.password = viper.GetString("database.password")
	databaseConfiguration.sslmode = viper.GetString("database.sslmode")
}

// Return the port of the database server
func (databaseConfiguration *databaseConfigurationImpl) GetPort() int {
	return databaseConfiguration.port
}

// Return the host of the database server
func (databaseConfiguration *databaseConfigurationImpl) GetHost() string {
	return databaseConfiguration.host
}

// Return the database name
func (databaseConfiguration *databaseConfigurationImpl) GetDBName() string {
	return databaseConfiguration.dbName
}

// Return the username of the user on the database server
func (databaseConfiguration *databaseConfigurationImpl) GetUsername() string {
	return databaseConfiguration.username
}

// Return the password of the user on the database server
func (databaseConfiguration *databaseConfigurationImpl) GetPassword() string {
	return databaseConfiguration.password
}

// Return the sslmode for the database connection
func (databaseConfiguration *databaseConfigurationImpl) GetSSLMode() string {
	return databaseConfiguration.sslmode
}

func (databaseConfiguration *databaseConfigurationImpl) Log() {
	zap.L().Info("Database configuration",
		zap.Int("port", databaseConfiguration.port),
		zap.String("host", databaseConfiguration.host),
		zap.String("dbName", databaseConfiguration.dbName),
		zap.String("username", databaseConfiguration.username),
		zap.String("password", databaseConfiguration.password),
		zap.String("sslmode", databaseConfiguration.sslmode),
	)
}
