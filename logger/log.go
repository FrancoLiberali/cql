package logger

import (
	"log"

	"go.uber.org/zap"

	"github.com/ditrit/badaas/configuration"
)

// Return a configured logger
func NewLogger(conf configuration.LoggerConfiguration) *zap.Logger {
	var config zap.Config
	if conf.GetMode() == configuration.ProductionLogger {
		config = zap.NewProductionConfig()

		log.Printf("Log mode use: %s\n", configuration.ProductionLogger)
	} else {
		config = zap.NewDevelopmentConfig()
		log.Printf("Log mode use: %s\n", configuration.DevelopmentLogger)
	}

	config.DisableStacktrace = conf.GetDisableStacktrace()

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	logger.Info("The logger was successfully initialized")

	return logger
}
