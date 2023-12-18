package orm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/model"
)

func GetCRUD[T model.Model, ID model.ID](db *gorm.DB) (CRUDService[T, ID], CRUDRepository[T, ID]) {
	repository := NewCRUDRepository[T, ID]()
	return NewCRUDService(db, repository), repository
}

func autoMigrate(modelsLists [][]any, db *gorm.DB) error {
	if len(modelsLists) > 0 {
		allModels := pie.Flat(modelsLists)
		return db.AutoMigrate(allModels...)
	}

	return nil
}
