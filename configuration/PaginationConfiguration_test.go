package configuration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
)

var PaginationConfigurationString = `server.pagination.page.max: 12`

func TestPaginationConfigurationNewPaginationConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewPaginationConfiguration(), "the contructor for PaginationConfiguration should not return a nil value")
}

func TestPaginationConfigurationGetMaxElemPerPage(t *testing.T) {
	setupViperEnvironment(PaginationConfigurationString)
	PaginationConfiguration := configuration.NewPaginationConfiguration()
	assert.Equal(t, uint(12), PaginationConfiguration.GetMaxElemPerPage())
}
