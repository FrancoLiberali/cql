package configuration

import "github.com/spf13/viper"

// Hold the configuration values for the database connection
type DatabaseConfiguration struct{}

// Instantiate a new configuration holder for the database connection
func NewDatabaseConfiguration() *DatabaseConfiguration {
	return &DatabaseConfiguration{}
}

// Return the port of the database server
func (dbc *DatabaseConfiguration) GetPort() int {
	return viper.GetInt("database.port")
}

// Return the host of the database server
func (dbc *DatabaseConfiguration) GetHost() string {
	return viper.GetString("database.host")
}

// Return the database name
func (dbc *DatabaseConfiguration) GetDBName() string {
	return viper.GetString("database.name")
}

// Return the username of the user on the database server
func (dbc *DatabaseConfiguration) GetUsername() string {
	return viper.GetString("database.username")
}

// Return the password of the user on the database server
func (dbc *DatabaseConfiguration) GetPassword() string {
	return viper.GetString("database.password")

}

// Return the sslmode for the database connection
func (dbc *DatabaseConfiguration) GetSSLMode() string {
	return viper.GetString("database.sslmode")
}
