package commands

import (
	"fmt"

	"github.com/ditrit/verdeter"
	verdetermodels "github.com/ditrit/verdeter/models"
)

func initLoggerCommands(cfg *verdeter.VerdeterCommand) {

	loggerModeKey := "logger.mode"
	cfg.GKey(loggerModeKey, verdeter.IsStr, "", "The logger mode (default to \"prod\")")
	cfg.SetDefault(loggerModeKey, "prod")
	cfg.SetRequired(loggerModeKey)
	// TODO: remove when this issue is released (https://github.com/ditrit/verdeter/issues/5)
	loggerModeValidator := verdetermodels.Validator{
		Name: "Logger Mode Validator",
		Func: func(input interface{}) error {
			allowedModes := []string{"prod", "dev"}
			inputStr, ok := input.(string)
			if !ok {
				return fmt.Errorf("logger.mode value should be a string")
			}
			if inputStr == "" {
				return fmt.Errorf("logger.mode value can't be empty")
			}
			for _, mode := range allowedModes {
				if mode == inputStr {
					return nil
				}
			}
			return fmt.Errorf("logger.mode should be one of the implemented modes (%v)", allowedModes)

		},
	}
	cfg.SetValidator(loggerModeKey, loggerModeValidator)

	loggerRequestTemplateKey := "logger.request.template"
	cfg.GKey(loggerRequestTemplateKey, verdeter.IsStr, "", "Template message for all request logs")
	cfg.SetDefault(loggerRequestTemplateKey, "Receive {{method}} request on {{url}}")
}
