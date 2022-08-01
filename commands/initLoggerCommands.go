package commands

import (
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/validators"
)

func initLoggerCommands(cfg *verdeter.VerdeterCommand) {
	cfg.GKey(configuration.LoggerModeKey, verdeter.IsStr, "", "The logger mode (default to \"prod\")")
	cfg.SetDefault(configuration.LoggerModeKey, "prod")
	cfg.AddValidator(configuration.LoggerModeKey, validators.AuthorizedValues("prod", "dev"))

	cfg.GKey(configuration.LoggerRequestTemplateKey, verdeter.IsStr, "", "Template message for all request logs")
	cfg.SetDefault(configuration.LoggerRequestTemplateKey, "Receive {{method}} request on {{url}}")
}
