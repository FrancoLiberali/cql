package orm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"
)

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	if len(modelsLists) > 0 {
		allModels := pie.Flat(modelsLists)
		return db.AutoMigrate(allModels...)
	}

	return nil
}
