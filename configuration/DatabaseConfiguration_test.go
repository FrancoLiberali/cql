package configuration_test

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/ditrit/badaas/configuration"
)

var databaseConfigurationString = `
database:
  host: e2e-db-1
  port: 26257
  sslmode: disable
  username: root
  password: postgres
  name: badaas_db
  init:
    retry: 10
    retryTime: 5
`

// Set the viper global instance config to the content of the string passed as argument
//
// First the content of viper global is wiped
// Then the content of the string (yml config) is parsed and read by viper
// For more information please head to: https://github.com/spf13/viper#reading-config-from-ioreader
func setupViperEnvironment(configurationString string) {
	viper.Reset()
	viper.SetConfigType("yaml")

	err := viper.ReadConfig(strings.NewReader(configurationString))
	if err != nil {
		panic(err)
	}
}

func TestDatabaseConfigurationNewDBConfig(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.NotNil(t, databaseConfiguration, "the database configuration should not be nil")
}

func TestDatabaseConfigurationGetPort(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, 26257, databaseConfiguration.GetPort(), "should be equals")
}

func TestDatabaseConfigurationGetHost(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, "e2e-db-1", databaseConfiguration.GetHost())
}

func TestDatabaseConfigurationGetUsername(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, "root", databaseConfiguration.GetUsername())
}

func TestDatabaseConfigurationGetPassword(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, "postgres", databaseConfiguration.GetPassword())
}

func TestDatabaseConfigurationGetSSLMode(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, "disable", databaseConfiguration.GetSSLMode())
}

func TestDatabaseConfigurationGetDBName(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, "badaas_db", databaseConfiguration.GetDBName())
}

func TestDatabaseConfigurationGetRetryTime(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, 5*time.Second, databaseConfiguration.GetRetryTime())
}

func TestDatabaseConfigurationGetRetry(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)

	databaseConfiguration := configuration.NewDatabaseConfiguration()
	assert.Equal(t, uint(10), databaseConfiguration.GetRetry())
}

func TestDatabaseConfigurationLog(t *testing.T) {
	setupViperEnvironment(databaseConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)
	databaseConfiguration := configuration.NewDatabaseConfiguration()
	databaseConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Database configuration", log.Message)
	require.Len(t, log.Context, 8)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "port", Type: zapcore.Int64Type, Integer: 26257},
		{Key: "retry", Type: zapcore.Uint64Type, Integer: 10},
		{Key: "retryTime", Type: zapcore.Uint64Type, Integer: 5},
		{Key: "host", Type: zapcore.StringType, String: "e2e-db-1"},
		{Key: "dbName", Type: zapcore.StringType, String: "badaas_db"},
		{Key: "username", Type: zapcore.StringType, String: "root"},
		{Key: "password", Type: zapcore.StringType, String: "postgres"},
		{Key: "sslmode", Type: zapcore.StringType, String: "disable"},
	}, log.Context)
}
