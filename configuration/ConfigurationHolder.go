package configuration

// Every configuration holder must implement this interface
type ConfigurationHolder interface {
	// Reload the values provided by the configuration holder
	Reload()

	// Log the values provided by the configuration holder using zap's global logger
	Log()
}
