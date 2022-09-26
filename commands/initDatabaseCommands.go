package commands

import "github.com/ditrit/verdeter"

func initDatabaseCommands(cfg *verdeter.VerdeterCommand) {
	databasePortKey := "database.port"
	cfg.GKey(databasePortKey, verdeter.IsInt, "", "The port of the database server")
	cfg.SetRequired(databasePortKey)

	databaseHostKey := "database.host"
	cfg.GKey(databaseHostKey, verdeter.IsStr, "", "The host of the database server")
	cfg.SetRequired(databaseHostKey)

	databaseNameKey := "database.name"
	cfg.GKey(databaseNameKey, verdeter.IsStr, "", "The name of the database to use")
	cfg.SetRequired(databaseNameKey)

	databaseUsernameKey := "database.username"
	cfg.GKey(databaseUsernameKey, verdeter.IsStr, "", "The username of the account on the database server")
	cfg.SetRequired(databaseUsernameKey)

	databasePasswordKey := "database.password"
	cfg.GKey(databasePasswordKey, verdeter.IsStr, "", "The password of the account one the database server")
	cfg.SetRequired(databasePasswordKey)

	databaseSslmodeKey := "database.sslmode"
	cfg.GKey(databaseSslmodeKey, verdeter.IsStr, "", "The sslmode to use when connecting to the database server")
	cfg.SetRequired(databaseSslmodeKey)
}
