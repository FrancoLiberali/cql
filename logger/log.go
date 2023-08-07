package logger

import (
	"log"

	"go.uber.org/zap"

	"github.com/ditrit/badaas/configuration"
)

const (
	ProductionLogger  = "prod"
	DevelopmentLogger = "dev"
)

// Return a configured logger
func NewLogger(conf configuration.LoggerConfiguration) *zap.Logger {
	var config zap.Config
	if conf.GetMode() == ProductionLogger {
		config = zap.NewProductionConfig()

		log.Printf("Log mode use: %s\n", ProductionLogger)
	} else {
		config = zap.NewDevelopmentConfig()

		log.Printf("Log mode use: %s\n", DevelopmentLogger)
	}

	config.DisableStacktrace = true

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	logger.Info("The logger was successfully initialized")

	return logger
}
