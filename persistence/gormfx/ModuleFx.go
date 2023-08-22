package gormfx

import (
	"github.com/elliotchance/pie/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
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

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	if len(modelsLists) > 0 {
		allModels := pie.Flat(modelsLists)
		return db.AutoMigrate(allModels...)
	}

	return nil
}
