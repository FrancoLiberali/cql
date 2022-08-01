package router

import (
	"github.com/ditrit/badaas/router/middlewares"
	"go.uber.org/fx"
)

// RouterModule for fx
var RouterModule = fx.Module(
	"router",
	// middlewares
	fx.Provide(middlewares.NewJSONController),
	fx.Provide(middlewares.NewMiddlewareLogger),

	fx.Provide(middlewares.NewAuthenticationMiddleware),

	// create router
	fx.Provide(SetupRouter),
)
