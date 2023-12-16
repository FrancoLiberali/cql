package configuration_test

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ditrit/badaas/configuration"
	configurationMocks "github.com/ditrit/badaas/mocks/configuration"
	"github.com/ditrit/verdeter"
)

var rootCommand = verdeter.BuildVerdeterCommand(verdeter.VerdeterConfig{
	Use:   "badaas",
	Short: "Backend and Distribution as a Service",
	Run:   doNothing,
})

func doNothing(_ *cobra.Command, _ []string) {}

func TestInitCommandsInitializerSetsAllKeysWithoutError(t *testing.T) {
	err := configuration.NewCommandInitializer(
		configuration.NewKeySetter(),
	).Init(rootCommand)
	assert.Nil(t, err)
}

func TestInitCommandsInitializerReturnsErrorWhenErrorOnKeySet(t *testing.T) {
	mockKeySetter := configurationMocks.NewKeySetter(t)
	mockKeySetter.On("Set", mock.Anything, mock.Anything).Return(errors.New("error setting key"))

	commandInitializer := configuration.NewCommandInitializer(mockKeySetter)

	err := commandInitializer.Init(rootCommand)
	assert.ErrorContains(t, err, "error setting key")
}
