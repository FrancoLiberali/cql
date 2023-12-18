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
		{
			Name:     LoggerDisableStacktraceKey,
			ValType:  verdeter.IsBool,
			Usage:    "Disable error stacktrace from logs (default to true)",
			DefaultV: true,
		},
		{
			Name:     LoggerSlowQueryThresholdKey,
			ValType:  verdeter.IsInt,
			Usage:    "Threshold for the slow query warning in milliseconds (default to 200)",
			DefaultV: defaultLoggerSlowQueryThreshold,
		},
		{
			Name:     LoggerSlowTransactionThresholdKey,
			ValType:  verdeter.IsInt,
			Usage:    "Threshold for the slow transaction warning in milliseconds (default to 200)",
			DefaultV: defaultLoggerSlowTransactionThreshold,
		},
		{
			Name:     LoggerIgnoreRecordNotFoundErrorKey,
			ValType:  verdeter.IsBool,
			Usage:    "If true, ignore gorm.ErrRecordNotFound error for logger (default to false)",
			DefaultV: false,
		},
		{
			Name:     LoggerParameterizedQueriesKey,
			ValType:  verdeter.IsBool,
			Usage:    "If true, don't include params in the query execution logs (default to false)",
			DefaultV: false,
		},
	}
}
