package configuration

import (
	"github.com/ditrit/verdeter"
)

// Definition of session configuration keys
func getSessionConfigurationKeys() []KeyDefinition {
	return []KeyDefinition{
		{
			Name:     SessionDurationKey,
			ValType:  verdeter.IsUint,
			Usage:    "The duration of a user session in seconds",
			DefaultV: defaultSessionDuration,
		},
		{
			Name:     SessionPullIntervalKey,
			ValType:  verdeter.IsUint,
			Usage:    "The refresh interval in seconds. Badaas refresh it's internal session cache periodically",
			DefaultV: defaultSessionPullInterval,
		},
		{
			Name:     SessionRollIntervalKey,
			ValType:  verdeter.IsUint,
			Usage:    "The interval in which the user can renew it's session by making a request",
			DefaultV: defaultSessionRollInterval,
		},
	}
}
