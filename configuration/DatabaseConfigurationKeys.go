package configuration

import (
	"github.com/ditrit/verdeter"
)

// Definition of database configuration keys
func getDatabaseConfigurationKeys() []KeyDefinition {
	return []KeyDefinition{
		{
			Name:     DatabasePortKey,
			ValType:  verdeter.IsInt,
			Usage:    "The port of the database server",
			Required: true,
		},
		{
			Name:     DatabaseHostKey,
			ValType:  verdeter.IsStr,
			Usage:    "The host of the database server",
			Required: true,
		},
		{
			Name:    DatabaseNameKey,
			ValType: verdeter.IsStr,
			Usage:   "The name of the database to use",
		},
		{
			Name:     DatabaseSslmodeKey,
			ValType:  verdeter.IsStr,
			Usage:    "The sslmode to use when connecting to the database server",
			Required: true,
		},
		{
			Name:     DatabaseUsernameKey,
			ValType:  verdeter.IsStr,
			Usage:    "The username of the account on the database server",
			Required: true,
		},
		{
			Name:     DatabasePasswordKey,
			ValType:  verdeter.IsStr,
			Usage:    "The password of the account one the database server",
			Required: true,
		},
		{
			Name:     DatabaseRetryKey,
			ValType:  verdeter.IsUint,
			Usage:    "The number of times badaas tries to establish a connection with the database",
			DefaultV: uint(10),
		},
		{
			Name:     DatabaseRetryDurationKey,
			ValType:  verdeter.IsUint,
			Usage:    "The duration in seconds badaas wait between connection attempts",
			DefaultV: uint(5),
		},
	}
}
