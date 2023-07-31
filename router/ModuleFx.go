package router

import (
	"go.uber.org/fx"

	"github.com/ditrit/badaas/router/middlewares"
)

// RouterModule for fx
var RouterModule = fx.Module(
	"router",
	fx.Provide(NewRouter),
	// middlewares
	fx.Provide(middlewares.NewJSONController),
	fx.Provide(middlewares.NewMiddlewareLogger),
	fx.Invoke(middlewares.AddLoggerMiddleware),
)
