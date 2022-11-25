package configuration

import "go.uber.org/fx"

// ConfigurationModule for fx
var ConfigurationModule = fx.Module(
	"configuration",
	fx.Provide(NewDatabaseConfiguration),
	fx.Provide(NewHTTPServerConfiguration),
	fx.Provide(NewLoggerConfiguration),
	fx.Provide(NewPaginationConfiguration),
	fx.Provide(NewInitializationConfiguration),
)
