package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	configurationmocks "github.com/ditrit/badaas/mocks/configuration"
)

func TestInitializeDevelopmentLogger(t *testing.T) {
	conf := configurationmocks.NewLoggerConfiguration(t)
	conf.On("GetMode").Return("dev")
	logger := NewLogger(conf)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))
}

func TestInitializeProductionLogger(t *testing.T) {
	conf := configurationmocks.NewLoggerConfiguration(t)
	conf.On("GetMode").Return("prod")
	logger := NewLogger(conf)
	assert.NotNil(t, logger)
	assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))
}

func TestInitializeProductionLoggerNoConf(t *testing.T) {
	conf := configurationmocks.NewLoggerConfiguration(t)
	conf.On("GetMode").Return("a stupid value")
	logger := NewLogger(conf)
	assert.NotNil(t, logger)
	assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))
}
