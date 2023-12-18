package configuration

import "go.uber.org/zap"

// Every configuration holder must implement this interface
type Holder interface {
	// Reload the values provided by the configuration holder
	Reload()

	// Log the values provided by the configuration holder
	Log(logger *zap.Logger)
}
