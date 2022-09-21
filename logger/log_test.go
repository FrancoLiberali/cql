package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitializeDevelopmentLogger(t *testing.T) {
	err := initLogger(DevelopmentLogger)
	if err != nil {
		t.Errorf("InitLogger should return a null error")
	}
	assert.True(t, zap.L().Core().Enabled(zapcore.DebugLevel))
}

func TestInitializeProductionLogger(t *testing.T) {
	err := initLogger(ProductionLogger)
	if err != nil {
		t.Errorf("InitLogger should return a null error")
	}
	assert.False(t, zap.L().Core().Enabled(zapcore.DebugLevel))

}
