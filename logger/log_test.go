package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/ditrit/badaas/configuration"
	configurationmocks "github.com/ditrit/badaas/mocks/configuration"
)

func TestInitializeDevelopmentLogger(t *testing.T) {
	conf := configurationmocks.NewLoggerConfiguration(t)
	conf.On("GetMode").Return(configuration.DevelopmentLogger)
	conf.On("GetDisableStacktrace").Return(true)
	logger := NewLogger(conf)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))
}

func TestInitializeProductionLogger(t *testing.T) {
	conf := configurationmocks.NewLoggerConfiguration(t)
	conf.On("GetMode").Return(configuration.ProductionLogger)
	conf.On("GetDisableStacktrace").Return(true)
	logger := NewLogger(conf)
	assert.NotNil(t, logger)
	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))
}

func TestInitializeProductionLoggerNoConf(t *testing.T) {
	conf := configurationmocks.NewLoggerConfiguration(t)
	conf.On("GetMode").Return("a stupid value")
	conf.On("GetDisableStacktrace").Return(true)
	logger := NewLogger(conf)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))
}
