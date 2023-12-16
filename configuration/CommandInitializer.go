package configuration

import (
	"github.com/ditrit/verdeter"
)

type CommandInitializer interface {
	// Inits VerdeterCommand "command" with the all the keys that are configurable in badaas
	Init(command *verdeter.VerdeterCommand) error
}

// Implementation of the CommandInitializer
type commandInitializerImpl struct {
	KeySetter KeySetter
	Keys      []KeyDefinition
}

// Constructor of CommandInitializer with the keys for badaas
// it uses the keySetter to set the configuration keys in the VerdeterCommand
func NewCommandInitializer(keySetter KeySetter) CommandInitializer {
	keys := []KeyDefinition{
		{
			Name:     "config_path",
			ValType:  verdeter.IsStr,
			Usage:    "Path to the config file/directory",
			DefaultV: ".",
		},
	}
	keys = append(keys, getDatabaseConfigurationKeys()...)
	keys = append(keys, getSessionConfigurationKeys()...)
	keys = append(keys, getInitializationConfigurationKeys()...)
	keys = append(keys, getServerConfigurationKeys()...)
	keys = append(keys, getLoggerConfigurationKeys()...)

	return commandInitializerImpl{
		KeySetter: keySetter,
		Keys:      keys,
	}
}

// Inits VerdeterCommand "cmd" with the all the keys in the Keys of the initializer
func (initializer commandInitializerImpl) Init(command *verdeter.VerdeterCommand) error {
	for _, key := range initializer.Keys {
		err := initializer.KeySetter.Set(command, key)
		if err != nil {
			return err
		}
	}

	return nil
}
