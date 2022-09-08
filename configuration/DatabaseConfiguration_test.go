package configuration_test

import (
	"strings"
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var databaseConfigurationString = `
database:
  host: e2e-db-1
  port: 26257
  sslmode: disable
  username: root
  password: postgres
  name: badaas_db
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
