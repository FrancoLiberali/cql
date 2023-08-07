package configuration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/ditrit/badaas/configuration"
)

var PaginationConfigurationString = `server.pagination.page.max: 12`

func TestPaginationConfigurationNewPaginationConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewPaginationConfiguration(), "the constructor for PaginationConfiguration should not return a nil value")
}

func TestPaginationConfigurationGetMaxElemPerPage(t *testing.T) {
	setupViperEnvironment(PaginationConfigurationString)

	paginationConfiguration := configuration.NewPaginationConfiguration()
	assert.Equal(t, uint(12), paginationConfiguration.GetMaxElemPerPage())
}

func TestPaginationConfigurationLog(t *testing.T) {
	setupViperEnvironment(PaginationConfigurationString)
	// creating logger
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	observedLogger := zap.New(observedZapCore)

	paginationConfiguration := configuration.NewPaginationConfiguration()
	paginationConfiguration.Log(observedLogger)

	require.Equal(t, 1, observedLogs.Len())
	log := observedLogs.All()[0]
	assert.Equal(t, "Pagination configuration", log.Message)
	require.Len(t, log.Context, 1)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "maxelemPerPage", Type: zapcore.Uint64Type, Integer: 12},
	}, log.Context)
}
