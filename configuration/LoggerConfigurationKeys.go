package configuration

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

// Definition of logger configuration keys
func getLoggerConfigurationKeys() []KeyDefinition {
	modeValidator := validators.AuthorizedValues("prod", "dev")
	return []KeyDefinition{
		{
			Name:     LoggerRequestTemplateKey,
			ValType:  verdeter.IsStr,
			Usage:    "Template message for all request logs",
			DefaultV: "Receive {{method}} request on {{url}}",
		},
		{
			Name:      LoggerModeKey,
			ValType:   verdeter.IsStr,
			Usage:     "The logger mode (default to \"prod\")",
			DefaultV:  "prod",
			Validator: &modeValidator,
		},
	}
}
