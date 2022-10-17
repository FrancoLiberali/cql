package logger

import (
	"log"

	"github.com/ditrit/badaas/configuration"
	"go.uber.org/zap"
)

const (
	ProductionLogger  = "prod"
	DevelopmentLogger = "dev"
)

// Initialize zap global logger instance
func initLogger(mode string) error {
	var config zap.Config
	var err error
	if mode == ProductionLogger {
		config = zap.NewProductionConfig()
		log.Printf("Log mode use: %s\n", ProductionLogger)

	} else {
		config = zap.NewDevelopmentConfig()
		log.Printf("Log mode use: %s\n", DevelopmentLogger)

	}
	config.DisableStacktrace = true
	logger, err := config.Build()
	if err == nil {
		zap.ReplaceGlobals(logger)
	}
	return err
}

// Initialize zap global logger instance
func InitLoggerFromConf() error {
	loggerConfiguration := configuration.Get().LoggerConfiguration
	return initLogger(loggerConfiguration.GetMode())
}
