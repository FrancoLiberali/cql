package orm

import (
	"go.uber.org/fx"
)

type GetModelsResult struct {
	fx.Out

	Models []any `group:"modelsTables"`
}

var AutoMigrate = fx.Module(
	"AutoMigrate",
	fx.Invoke(
		fx.Annotate(
			autoMigrate,
			fx.ParamTags(`group:"modelsTables"`),
		),
	),
)
