package configuration

// The Configuration instance that is shared across the codebase (via the Get() function)
var localConfiguration *Configuration

// Configuration struct
//
// Hold all the differents configuration holders
type Configuration struct {
	ConfigurationHolder

	HTTPServerConfiguration HTTPServerConfiguration
	LoggerConfiguration     LoggerConfiguration
	DatabaseConfiguration   DatabaseConfiguration
	PaginationConfiguration PaginationConfiguration
}

// Configuration constructor
//
// Create all configuration holder and reload them
func NewConfiguration() *Configuration {
	return &Configuration{
		HTTPServerConfiguration: NewHTTPServerConfiguration(),
		LoggerConfiguration:     NewLoggerConfiguration(),
		DatabaseConfiguration:   NewDatabaseConfiguration(),
		PaginationConfiguration: NewPaginationConfiguration(),
	}

}

// Replace the global configuration instance
func ReplaceGlobals(configuration *Configuration) {
	localConfiguration = configuration
}

// Get global instance
func Get() *Configuration {
	if localConfiguration == nil {
		panic("local configuration instance is nil")
	}
	return localConfiguration
}

// Log all configuration configuration
func (configuration *Configuration) Log() {
	configuration.HTTPServerConfiguration.Log()
	configuration.LoggerConfiguration.Log()
	configuration.DatabaseConfiguration.Log()
	configuration.PaginationConfiguration.Log()
}

// Reload all configuration
func (configuration *Configuration) Reload() {
	configuration.HTTPServerConfiguration.Reload()
	configuration.LoggerConfiguration.Reload()
	configuration.DatabaseConfiguration.Reload()
	configuration.PaginationConfiguration.Reload()
}
