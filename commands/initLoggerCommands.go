package commands

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func initLoggerCommands(cfg *verdeter.VerdeterCommand) {

	loggerModeKey := "logger.mode"
	cfg.GKey(loggerModeKey, verdeter.IsStr, "", "The logger mode (default to \"prod\")")
	cfg.SetDefault(loggerModeKey, "prod")
	cfg.AddValidator(loggerModeKey, validators.AuthorizedValues("authorized values", "prod", "dev"))

	loggerRequestTemplateKey := "logger.request.template"
	cfg.GKey(loggerRequestTemplateKey, verdeter.IsStr, "", "Template message for all request logs")
	cfg.SetDefault(loggerRequestTemplateKey, "Receive {{method}} request on {{url}}")
}
