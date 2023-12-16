package configuration

import (
	"time"

	"github.com/ditrit/badaas/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// The config keys regarding the database settings
const (
	DatabasePortKey          string = "database.port"
	DatabaseHostKey          string = "database.host"
	DatabaseNameKey          string = "database.name"
	DatabaseUsernameKey      string = "database.username"
	DatabasePasswordKey      string = "database.password"
	DatabaseSslmodeKey       string = "database.sslmode"
	DatabaseRetryKey         string = "database.init.retry"
	DatabaseRetryDurationKey string = "database.init.retryTime"
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
	GetRetry() uint
	GetRetryTime() time.Duration
}

// Concrete implementation of the DatabaseConfiguration interface
type databaseConfigurationImpl struct {
	port      int
	host      string
	dbName    string
	username  string
	password  string
	sslmode   string
	retry     uint
	retryTime uint
}

// Instantiate a new configuration holder for the database connection
func NewDatabaseConfiguration() DatabaseConfiguration {
	databaseConfiguration := new(databaseConfigurationImpl)
	databaseConfiguration.Reload()
	return databaseConfiguration
}

// Reload database configuration
func (databaseConfiguration *databaseConfigurationImpl) Reload() {
	databaseConfiguration.port = viper.GetInt(DatabasePortKey)
	databaseConfiguration.host = viper.GetString(DatabaseHostKey)
	databaseConfiguration.dbName = viper.GetString(DatabaseNameKey)
	databaseConfiguration.username = viper.GetString(DatabaseUsernameKey)
	databaseConfiguration.password = viper.GetString(DatabasePasswordKey)
	databaseConfiguration.sslmode = viper.GetString(DatabaseSslmodeKey)
	databaseConfiguration.retry = viper.GetUint(DatabaseRetryKey)
	databaseConfiguration.retryTime = viper.GetUint(DatabaseRetryDurationKey)
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

// Return the number of retries for the database connection
func (databaseConfiguration *databaseConfigurationImpl) GetRetry() uint {
	return databaseConfiguration.retry
}

// Return the waiting time between the database connections in seconds
func (databaseConfiguration *databaseConfigurationImpl) GetRetryTime() time.Duration {
	return utils.IntToSecond(int(databaseConfiguration.retryTime))
}

// Log the values provided by the configuration holder
func (databaseConfiguration *databaseConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("Database configuration",
		zap.Int("port", databaseConfiguration.port),
		zap.String("host", databaseConfiguration.host),
		zap.String("dbName", databaseConfiguration.dbName),
		zap.String("username", databaseConfiguration.username),
		zap.String("password", databaseConfiguration.password),
		zap.String("sslmode", databaseConfiguration.sslmode),
		zap.Uint("retry", databaseConfiguration.retry),
		zap.Uint("retryTime", databaseConfiguration.retryTime),
	)
}
