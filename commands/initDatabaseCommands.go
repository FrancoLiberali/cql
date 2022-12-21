package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
)

func initDatabaseCommands(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.DatabasePortKey, verdeter.IsInt, "", "The port of the database server")
	cfg.SetRequired(configuration.DatabasePortKey)

	cfg.GKey(configuration.DatabaseHostKey, verdeter.IsStr, "", "The host of the database server")
	cfg.SetRequired(configuration.DatabaseHostKey)

	cfg.GKey(configuration.DatabaseNameKey, verdeter.IsStr, "", "The name of the database to use")
	cfg.SetRequired(configuration.DatabaseNameKey)

	cfg.GKey(configuration.DatabaseUsernameKey, verdeter.IsStr, "", "The username of the account on the database server")
	cfg.SetRequired(configuration.DatabaseUsernameKey)

	cfg.GKey(configuration.DatabasePasswordKey, verdeter.IsStr, "", "The password of the account one the database server")
	cfg.SetRequired(configuration.DatabasePasswordKey)

	cfg.GKey(configuration.DatabaseSslmodeKey, verdeter.IsStr, "", "The sslmode to use when connecting to the database server")
	cfg.SetRequired(configuration.DatabaseSslmodeKey)

	cfg.GKey(configuration.DatabaseRetryKey, verdeter.IsUint, "", "The number of times badaas tries to establish a connection with the database")
	cfg.SetDefault(configuration.DatabaseRetryKey, uint(10))

	cfg.GKey(configuration.DatabaseRetryDurationKey, verdeter.IsUint, "", "The duration in seconds badaas wait between connection attempts")
	cfg.SetDefault(configuration.DatabaseRetryDurationKey, uint(5))
}
