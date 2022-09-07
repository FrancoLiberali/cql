package configuration_test

import (
	"testing"

	"github.com/ditrit/badaas/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpServerConfiguration(t *testing.T) {
	assert.NotNil(t, configuration.NewHTTPServerConfiguration(), "the contructor for HttpServerConfiguration should not return a nil value")
}
