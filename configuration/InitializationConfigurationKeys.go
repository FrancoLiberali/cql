package configuration

import (
	"github.com/ditrit/verdeter"
)

// Definition of initialization configuration keys
func getInitializationConfigurationKeys() []KeyDefinition {
	return []KeyDefinition{
		{
			Name:     InitializationDefaultAdminPasswordKey,
			ValType:  verdeter.IsStr,
			Usage:    "Set the default admin password is the admin user is not created yet.",
			DefaultV: "admin",
		},
	}
}
