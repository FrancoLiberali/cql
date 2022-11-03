package logger

import "go.uber.org/fx"

// LoggerModule for fx
var LoggerModule = fx.Module(
	"logger",
	fx.Provide(NewLogger),
)
