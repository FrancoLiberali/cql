package configuration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	conf := configuration.NewConfiguration()
	assert.NotNil(t, conf)
	// assert that the configurations are well instanciated
	assert.NotPanics(t, func() { conf.DatabaseConfiguration.Log() })
	assert.NotPanics(t, func() { conf.LoggerConfiguration.Log() })
	assert.NotPanics(t, func() { conf.HTTPServerConfiguration.Log() })
	assert.NotPanics(t, func() { conf.PaginationConfiguration.Log() })
}
