package configuration

import (
	"github.com/ditrit/verdeter"
	"github.com/ditrit/verdeter/models"
)

type KeySetter interface {
	// Configures the VerdeterCommand "command" with the information contained in "key"
	Set(command *verdeter.VerdeterCommand, key KeyDefinition) error
}

type KeyDefinition struct {
	Name      string
	ValType   models.ConfigType
	Usage     string
	Required  bool
	DefaultV  any
	Validator *models.Validator
}

type keySetterImpl struct{}

func NewKeySetter() KeySetter {
	return keySetterImpl{}
}

// Configures the VerdeterCommand "command" with the information contained in "key"
func (ks keySetterImpl) Set(command *verdeter.VerdeterCommand, key KeyDefinition) error {
	if err := command.GKey(key.Name, key.ValType, "", key.Usage); err != nil {
		return err
	}

	if key.Required {
		command.SetRequired(key.Name)
	}

	if key.DefaultV != nil {
		command.SetDefault(key.Name, key.DefaultV)
	}

	if key.Validator != nil {
		command.AddValidator(key.Name, *key.Validator)
	}

	return nil
}
